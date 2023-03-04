package worm

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
	name string

	id string

	path string

	file io.Writer

	status int

	Config config

	Chapters []CharpterHeadInfo

	CurrentIndex int

	BufferSize int

	charpterBuffer myBuffer

	ch chan bufferItem
}

func CreateTask(novelName string, chapters []CharpterHeadInfo, config config) task {
	fmt.Println("新建task")
	goupsize := len(chapters) / 100

	if goupsize < 5 {
		goupsize = 5
	}

	file, err := os.OpenFile("./novel.txt", os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)
	if err != nil {
		panic("error in creat task :" + err.Error())
	}
	fmt.Println("bufferSize is ", goupsize)

	return task{
		Config:       config,
		Chapters:     chapters,
		CurrentIndex: -1,
		BufferSize:   goupsize,
		status:       statusStoped,
		ch:           make(chan bufferItem, goupsize),
		file:         file,
	}

}

func (t *task) Start() {
	fmt.Println("task开始启动")
	for i := 0; i < t.BufferSize; i++ {
		if i <= len(t.Chapters) {
			go t.getOneChapter(i, t.Chapters[i].CharpterURL)

		}

	}

	for {
		if t.CurrentIndex >= len(t.Chapters) {
			return
		}
		var item = <-t.ch

		fmt.Println("从通道获得:", item.index)

		t.charpterBuffer.inputcontent(item)

		if item.index-t.CurrentIndex == 1 {
			ok := t.storeFromBuffer(item.index)
			if ok {
				index := t.CurrentIndex + 1
				go t.getOneChapter(index, t.Chapters[index].CharpterURL)
			}
		}
	}

}

func (t *task) storeFromBuffer(startfrom int) bool {

	var data []byte

	var changed = false

	for {

		var content *string

		i, content := t.charpterBuffer.getContent(startfrom)

		if content == nil {

			if changed {
				return true
			}
			return false
		}

		t.file.Write([]byte("第" + fmt.Sprint(startfrom) + "章" + t.Chapters[startfrom].ChapterName + "\n"))

		data = []byte(*content)

		t.file.Write(data)

		t.file.Write([]byte("\n"))

		t.CurrentIndex += 1

		fmt.Println(t.Chapters[startfrom].ChapterName + " 已录入")

		t.charpterBuffer.deleteItem(i)

		changed = true

		startfrom += 1

	}
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
		content = s.Text()
	})

	fmt.Println("第", intdex, "章，总长：", len(content))

	t.ch <- bufferItem{index: intdex, content: &content}

	return nil
}
