package worm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/PuerkitoBio/goquery"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tidwall/gjson"
)

const (
	statusStoped int = iota
	statusRunning
	statusFinished
)

type bufferItem struct {
	index   int
	content *string
}

var nullItem = bufferItem{
	index:   -1,
	content: nil,
}

type myBuffer struct {
	buffer []bufferItem
}

func (b *myBuffer) getContent(index int) (int, *string) {
	for i, v := range b.buffer {
		if v.index == index {
			return i, v.content
		}
	}
	return -1, nil
}

func (b *myBuffer) inputcontent(item bufferItem) {
	for i, v := range b.buffer {
		if v == nullItem {
			b.buffer[i] = item
		}
	}

	b.buffer = append(b.buffer, item)

}

func (b *myBuffer) deleteItem(indexInBuffer int) {
	b.buffer[indexInBuffer] = nullItem
}

type task struct {
	Name         string `json:"name"`
	Id           string `json:"id"`
	CurrentIndex int    `json:"currentIndex"`
	BufferSize   int    `json:"bufferSize"`
	Status       int    `json:"status"`

	toClient       chan InputData
	path           string
	file           *os.File
	mark           string
	config         config
	chapters       []CharpterHeadInfo
	charpterBuffer myBuffer
	ch             chan bufferItem
}

var tasks []*task

func init() {

	// a, filepath, c, d := runtime.Caller(0)
	// fmt.Printf("a is %v b is %v c is %v d si %v", a, filepath, c, d)

	// dir, _ := path.Split(filepath)

	// for path.Base(dir) != "videosite" {
	// 	dir = dir[:len(dir)-1]
	// 	fmt.Println(dir)
	// 	dir, _ = path.Split(dir)
	// 	fmt.Println(dir)
	// }
	// novelPath := path.Join(dir, "novel")
	novelPath := "./novel"

	novels, err := ioutil.ReadDir(novelPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range novels {
		markPath := path.Join(novelPath, v.Name(), "mark.json")
		chaptersPath := path.Join(novelPath, v.Name(), "chapters.json")

		markFile, err := os.Open(markPath)
		if err != nil {
			fmt.Println(err)
		}
		markBytes, err := io.ReadAll(markFile)
		if err != nil {
			fmt.Println(err)
		}

		var t = *new(task)
		err = json.Unmarshal(markBytes, &t)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("task path is:", t.path)

		//读取章节信息
		chaptersFile, err := os.Open(chaptersPath)
		if err != nil {
			fmt.Println(err)
		}

		chaptersBytes, err := io.ReadAll(chaptersFile)
		if err != nil {
			fmt.Println(err)
		}

		var chapters []CharpterHeadInfo
		err = json.Unmarshal(chaptersBytes, &chapters)
		if err != nil {
			fmt.Println(err)
		}

		t.chapters = chapters

		t.toClient = make(chan InputData)

		tasks = append(tasks, &t)
		go t.sendToClient()

		fmt.Println("Tasks length is ", len(tasks))

	}

	go CollectTask()

}

func CreateTask(novelName string, chapters []CharpterHeadInfo, config config) (task, error) {

	fmt.Println("新建task")
	goupsize := len(chapters) / 100

	seconds := time.Now().Unix()

	folderName := novelName + "_" + fmt.Sprint(seconds)

	folderPath := path.Join("novel", folderName)

	if goupsize < 5 {
		goupsize = 5
	}

	_, callerPath, _, _ := runtime.Caller(0)
	dir, _ := path.Split(callerPath)

	for path.Base(dir) != "videosite" {
		dir = dir[:len(dir)-1]
		fmt.Println(dir)
		dir, _ = path.Split(dir)
		fmt.Println(dir)
	}

	folderPath = path.Join(dir, folderPath)

	err := os.MkdirAll(folderPath, 0750)
	if err != nil {
		return task{}, err
	}

	filePath := path.Join(folderPath, "novel.txt")
	markPath := path.Join(folderPath, "mark.json")

	// filePath = path.Join(dir, filePath)
	// markPath = path.Join(dir, markPath)

	os.Chdir(folderPath)

	_, err = os.Create("novel.txt")
	if err != nil {
		return task{}, err
	}

	mark, err := os.Create("mark.json")
	if err != nil {
		return task{}, err
	}

	chaptersFile, err := os.Create("chapters.json")
	if err != nil {
		return task{}, err
	}

	fmt.Println("bufferSize is ", goupsize)

	id, _ := gonanoid.New()
	var t = task{
		Id:           id,
		Name:         novelName,
		config:       config,
		chapters:     chapters,
		CurrentIndex: -1,
		BufferSize:   goupsize,
		Status:       statusStoped,
		ch:           make(chan bufferItem, goupsize),
		mark:         markPath,
		path:         filePath,
		toClient:     make(chan InputData),
	}

	tasks = append(tasks, &t)
	go t.sendToClient()

	data, err := json.Marshal(&t)
	if err != nil {
		fmt.Println(err)
	}
	_, err = mark.Write(data)
	if err != nil {
		fmt.Println("init mark file ", err)
	}
	mark.Close()

	data, err = json.Marshal(&chapters)
	if err != nil {
		fmt.Println(err)
	}

	_, err = chaptersFile.Write(data)
	if err != nil {
		fmt.Println("write chapterInfo in chapters", err)
	}

	chaptersFile.Close()

	return t, nil
}

func (t *task) Start() {
	if t.Status == statusFinished {
		fmt.Printf("novel '%v' is already finished", t.Name)
		return
	}

	t.Status = statusRunning

	novelFile, err := os.OpenFile(t.path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	t.charpterBuffer = *new(myBuffer)

	t.file = novelFile

	var gettingIndex = -1
	fmt.Println("task开始启动")
	for i := 0; i < t.BufferSize; i++ {
		if i <= len(t.chapters) {
			gettingIndex += 1

			go t.getOneChapter(i, t.chapters[gettingIndex].CharpterURL)

		}

	}

	for {

		var item = <-t.ch

		fmt.Println("从通道获得:", item.index)

		t.charpterBuffer.inputcontent(item)

		if item.index-t.CurrentIndex == 1 {
			savedum := t.storeFromBuffer(item.index)
			for i := 0; i < savedum; i++ {
				t.CurrentIndex += 1

				if gettingIndex < len(t.chapters)-1 {
					gettingIndex += 1

					go t.getOneChapter(gettingIndex, t.chapters[gettingIndex].CharpterURL)
				}

			}

			t.sendToClient()

			if t.CurrentIndex >= len(t.chapters)-1 {
				t.Status = statusFinished
				t.file.Close()

				t.markFinish()
				fmt.Println("爬取结束")
				return
			}

		}
	}

}

func (t *task) sendToClient() {

	data, err := json.Marshal(*t)
	if err != nil {
		fmt.Println(err)
	}

	jsonString := gjson.ParseBytes(data).String()

	out := InputData{Id: t.Id, Json: jsonString}

	fmt.Println("in sendToClient write into chan:", jsonString)

	t.toClient <- out
}

func (t *task) storeFromBuffer(startfrom int) int {

	var savedNum = 0

	var data []byte

	for {

		var content *string

		i, content := t.charpterBuffer.getContent(startfrom)

		if content == nil {

			return savedNum
		}

		_, err := t.file.Write([]byte("第" + fmt.Sprint(startfrom) + "章" + t.chapters[startfrom].ChapterName + "\n"))
		if err != nil {
			fmt.Println("write chapter header", err)
		}
		data = []byte(*content)

		_, err = t.file.Write(data)
		if err != nil {
			fmt.Println("write chapter data ", err)
		}
		_, err = t.file.Write([]byte("\n"))
		if err != nil {
			fmt.Println("write return", err)
		}

		t.markIndex(startfrom)

		fmt.Println(t.chapters[startfrom].ChapterName + " 已录入")

		t.charpterBuffer.deleteItem(i)

		startfrom += 1

		savedNum += 1

		if startfrom >= len(t.chapters) {
			t.markIndex(-2)
			t.Status = statusFinished
			return savedNum
		}

	}
}

func (t *task) markIndex(index int) bool {

	markFile, err := os.OpenFile(t.mark, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open mark file ", err)
		return false
	}
	defer markFile.Close()

	data, err := json.Marshal(t)
	if err != nil {
		fmt.Println("json mark ", err)
		return false
	}
	_, err = markFile.Write(data)
	if err != nil {
		fmt.Println("write mark ", err)
		return false
	}

	return true
}

func (t *task) getOneChapter(intdex int, url string) error {

	charpter1 := url //"https://www.dengbi.com/16/16686/11148312.html"这是第一章
	var content string

	fmt.Println("开始获取：" + url)

	res, err := http.Get(t.config.WebsetURl + charpter1)
	if err != nil {
		fmt.Println(err.Error())
	}

	body := res.Body
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var selector *goquery.Selection = doc.Selection

	for _, s := range t.config.ContentSelector {
		selector = selector.Find(s)
	}

	selector.Each(func(i int, s *goquery.Selection) {
		elementContent := s.Text()
		content += elementContent + "\n"
	})

	fmt.Println("第", intdex, "章，总长：", len(content))

	t.ch <- bufferItem{index: intdex, content: &content}

	return nil
}

func (t *task) markFinish() {

	markFile, err := os.OpenFile(t.mark, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open mark file ", err)
	}

	data, err := json.Marshal(t)
	if err != nil {
		fmt.Println("json mark ", err)
	}
	_, err = markFile.Write(data)
	if err != nil {
		fmt.Println("write mark ", err)
	}
	markFile.Close()

}
