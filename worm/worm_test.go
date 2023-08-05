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

// func TestGetDeWaterMark(t *testing.T) {
// 	var testText = "1111111 1111111 11111 1111如果被/浏/览/器/强/制进入它们的阅/读/模/式了,阅读体/验极/差请退出转/码阅读."
// 	config, err := getConfigObj("https://www.dengbi8.com/shu/134780/")
// 	if err != nil {
// 		t.Error("config不存在")
// 	}

// }
