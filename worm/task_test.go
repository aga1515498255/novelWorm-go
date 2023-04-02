package worm

import "testing"

func TestCreateTask(t *testing.T) {

	var fakeChapters []CharpterHeadInfo
	fakeChapters = append(fakeChapters, CharpterHeadInfo{ChapterName: "testName", CharpterURL: "test.html"})

	var tasklen = len(OutputTasks)
	_, err := CreateTask("test", fakeChapters, configs[0])

	if err != nil {
		t.Error(err)
	}

	var newlen = len(OutputTasks)

	if newlen-tasklen != 1 {
		t.Error("创建失败")
	}
}

// func TestTaskIinit(t *testing.T) {

// }
