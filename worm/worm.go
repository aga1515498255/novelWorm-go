package worm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	MODE_PREVIWE = 0
	MODE_FULL    = 1

	MARK_MID = 0
	MARK_END = 1
)

type WaterMark struct {
	Content string `json:"content"`
	Sort    int    `json:"sort"`
}

var configs []config

type config struct {
	Name            string      `json:"name"`
	WebsetURl       string      `json:"websetURl"`
	WaterMark       []WaterMark `json:"waterMark"`
	ChapterSelector []string    `json:"chapterSelector"`
	ContentSelector []string    `json:"contentSelector"`
	NameSelector    []string    `json:"nameSelector"`

	URLselector struct {
		ChapterName string `json:"chapterName"`
		ChapterURL  string `json:"chapterURL"`
	} `json:"urlSelector"`
}

func (c *config) getName(doc *goquery.Document) string {

	var selector *goquery.Selection = doc.Selection

	for _, s := range c.NameSelector {
		selector = selector.Find(s)
	}

	name := selector.Text()
	return name
}

func (c *config) getchapterRef(doc *goquery.Document) []CharpterHeadInfo {

	var chapterRef []CharpterHeadInfo
	var selector *goquery.Selection = doc.Selection

	for _, s := range c.ChapterSelector {
		selector = selector.Find(s)
	}

	selector.Each(func(i int, s *goquery.Selection) {
		charpterName := ""
		var CPInfo CharpterHeadInfo = CharpterHeadInfo{}
		ref, _ := s.Attr(c.URLselector.ChapterURL)

		if c.URLselector.ChapterName == "text" {

			charpterName = s.Text()
		}

		CPInfo.ChapterName = charpterName
		CPInfo.CharpterURL = ref

		chapterRef = append(chapterRef, CPInfo)
	})

	fmt.Println("chapter number is ", len(chapterRef))

	return chapterRef
}

func GetConfigs() []config {
	return configs
}

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

	dir := "."

	configPath := path.Join(dir, "config.json")

	fmt.Println("dir is ", dir)

	configfile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}

	configdata, err := ioutil.ReadAll(configfile)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(configdata, &configs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(configs)

}

type CharpterHeadInfo struct {
	ChapterName string
	CharpterURL string
}

func GetName(url string) string {
	config, err := getConfigObj(url)
	if err != nil {
		fmt.Println(err)
	}

	targetPage, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

	}

	body := targetPage.Body

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)

	}
	name := config.getName(doc)

	return name
}

func GetNovel(url string) {
	config, err := getConfigObj(url)
	if err != nil {
		panic("error in get config: " + err.Error())

	}

	chapterInfos, err := GetchapterURL(url, MODE_FULL)
	if err != nil {
		panic("error in get GetchapterURL: " + err.Error())

	}

	name := GetName(url)

	fmt.Println("chapter infos is ", len(chapterInfos))

	task, _ := CreateTask(name, chapterInfos, config)

	go task.Start()
}

func getConfigObj(indexURL string) (config, error) {

	reg, err := regexp.Compile(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}`)
	if err != nil {
		fmt.Println("erro in compile regexp", err)
		return config{}, err
	}
	url := reg.FindString(indexURL)
	fmt.Println("target url is:" + url)
	for _, c := range configs {
		if c.WebsetURl == url || c.WebsetURl == url+"/" {
			return c, nil
		}
	}

	return config{}, errors.New("no website found")
}

func GetchapterURL(indexURL string, mode int) ([]CharpterHeadInfo, error) {

	configObj, err := getConfigObj(indexURL)
	if err != nil {

		error := errors.New("in get configObj" + err.Error())

		return nil, error
	}

	charpter1 := indexURL //"https://www.dengbi.com/16/16686/"

	targetPage, err := http.Get(charpter1)
	if err != nil {
		error := errors.New("in get H5 page of chapters" + err.Error())

		return nil, error
	}

	body := targetPage.Body

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		error := errors.New("in convert H5 into doc" + err.Error())

		return nil, error
	}

	fmt.Println("输出获取到的文档：", doc)

	var res []CharpterHeadInfo

	if mode == MODE_PREVIWE {
		res = configObj.getchapterRef(doc)[0:20]

	} else if mode == MODE_FULL {
		res = configObj.getchapterRef(doc)
	}

	return res, nil

}

func loudChapter(chapterURL string) (string, error) {
	charpter1 := chapterURL //"https://www.dengbi.com/16/16686/11148312.html"这是第一章
	var content string

	res, err := http.Get(charpter1)
	if err != nil {
		return "", err
	}

	body := res.Body
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	doc.Find("article#article").Each(func(i int, s *goquery.Selection) {
		content = s.Text()
	})

	return content, nil

}

func strDewaterMark(str string, waterMark []string) string {

	s := str
	for _, w := range waterMark {
		s = deWaterMark(s, w)
	}

	return s
}

func deWaterMark(str string, waterMark string) string {
	i := strings.Index(str, waterMark)

	len := len(waterMark)

	res := str[:i] + str[i+len:]

	return res
}

func (c *config) DeWaterMark(str string) string {
	for _, v := range c.WaterMark {
		fmt.Println("sort of waterMark is", v.Sort)
		if v.Sort == MARK_MID {

			return loopDeletStr(str, v.Content)
		} else {
			i := strings.Index(str, v.Content)
			if i == -1 {
				return str
			} else {
				return str[:i]
			}
		}
	}
	return str
}

func loopDeletStr(str string, target string) string {
	i := strings.Index(str, target)

	res := str

	for {
		if i == -1 {
			return res
		}
		len := len(target)

		res = str[:i] + str[i+len:]

		i = strings.Index(str, target)

	}
}
