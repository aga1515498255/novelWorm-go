package worm

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGetNovelName(t *testing.T) {

	config, err := getConfigObj("https://www.dengbi8.com/shu/134780/")
	if err != nil {
		t.Error("config不存在")
	}

	targetPage, err := http.Get("https://www.dengbi8.com/shu/134780/")
	if err != nil {

		t.Error(err)

	}

	body := targetPage.Body

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {

		t.Error(err)

	}
	name := config.getName(doc)

	if name != "夜的命名术" {
		t.Error("失败")
	}

	// var fakeChapters []CharpterHeadInfo
	// fakeChapters = append(fakeChapters, CharpterHeadInfo{ChapterName: "testName", CharpterURL: "test.html"})

	// var tasklen = len(Tasks)
	// _, err := CreateTask("test", fakeChapters, configs[0])

	// if err != nil {
	// 	t.Error(err)
	// }

	// var newlen = len(Tasks)

	// if newlen-tasklen != 1 {
	// 	t.Error("创建失败")
	// }
}
