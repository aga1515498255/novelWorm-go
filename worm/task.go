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
	Name string `json:"name"`

	Id string `json:"id"`

	Path string `json:"path"`

	file io.Writer

	Mark string `json:"mark"`

	Status int `json:"status"`

	Config config `json:"config"`

	chapters []CharpterHeadInfo

	CurrentIndex int `json:"currentIndex"`

	BufferSize int `json:"bufferSize"`

	charpterBuffer myBuffer

	ch chan bufferItem
}

var Tasks []task

func init() {

	a, filepath, c, d := runtime.Caller(0)
	fmt.Printf("a is %v b is %v c is %v d si %v", a, filepath, c, d)

	dir, _ := path.Split(filepath)

	for path.Base(dir) != "videosite" {
		dir = dir[:len(dir)-1]
		fmt.Println(dir)
		dir, _ = path.Split(dir)
		fmt.Println(dir)
	}
	novelPath := path.Join(dir, "novel")

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

		fmt.Println("task path is:", t.Path)

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

		Tasks = append(Tasks, t)

		fmt.Println("Tasks length is ", len(Tasks))

	}

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

	// file, err := os.OpenFile("./novel.txt", os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)
	// if err != nil {
	// 	panic("error in creat task :" + err.Error())
	// }
	// fmt.Println("bufferSize is ", goupsize)

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

	var t = task{
		Name:         novelName,
		Config:       config,
		chapters:     chapters,
		CurrentIndex: -1,
		BufferSize:   goupsize,
		Status:       statusStoped,
		ch:           make(chan bufferItem, goupsize),
		Mark:         markPath,
		Path:         filePath,
	}

	Tasks = append(Tasks, t)

	data, err := json.Marshal(&t)
	if err != nil {
		fmt.Println(err)
	}
	mark.Write(data)

	mark.Close()

	data, err = json.Marshal(&t)
	if err != nil {
		fmt.Println(err)
	}
	chaptersFile.Write(data)

	chaptersFile.Close()

	return t, nil
}

func (t *task) Start() {

	if t.Status == statusFinished {
		fmt.Printf("novel '%v' is already finished", t.Name)
		return
	}

	t.Status = statusRunning

	novelFile, err := os.Open(t.Path)
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

		if t.CurrentIndex >= len(t.chapters)-1 {
			fmt.Println("爬取结束")
			return
		}

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

		}
	}

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

		t.file.Write([]byte("第" + fmt.Sprint(startfrom) + "章" + t.chapters[startfrom].ChapterName + "\n"))

		data = []byte(*content)

		t.file.Write(data)

		t.file.Write([]byte("\n"))

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
	markFile, err := os.OpenFile(t.Mark, os.O_SYNC|os.O_TRUNC, os.ModeSocket)
	if err != nil {
		return false
	}

	if index == -2 {
		markFile.Write([]byte("done"))
		return true
	}

	data, err := json.Marshal(t)
	markFile.Write(data)

	markFile.Close()

	return true
}

func (t *task) getOneChapter(intdex int, url string) error {

	charpter1 := url //"https://www.dengbi.com/16/16686/11148312.html"这是第一章
	var content string

	fmt.Println("开始获取：" + url)

	res, err := http.Get(t.Config.WebsetURl + charpter1)
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

	for _, s := range t.Config.ContentSelector {
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
