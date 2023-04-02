package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	app "ui"

	worm "worm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Open(uri string) error {
	cmd := exec.Command("cmd", "/C", "start "+uri)
	return cmd.Run()
}

var e *echo.Echo

func main() {

	e = echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} \n",
	}))

	g := e.Group("/api")

	g.GET("/preview/:url", preview1)
	g.GET("/config", config1)
	g.GET("/tasks", getTask1)
	g.GET("/getNovel/*", getnovel)

	e.GET("/*", handlStatic)

	Open("http://localhost:4321")

	e.Logger.Debug(e.Start(":4321"))

}

var uiFS fs.FS

func init() {
	var err error
	uiFS, err = fs.Sub(app.UI, "build") //
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}
	fmt.Println(os.Args[0])
}

func handlStatic(c echo.Context) error {

	// fmt.Println("in handlStatic")

	path := c.Request().URL.Path

	if path == "/" || path == "\\" {
		path = "index.html"
	}

	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimPrefix(path, "/")

	fmt.Println("in handlStatic path=", path)

	file, err := uiFS.Open(path)
	if err != nil {

		file, _ = uiFS.Open("index.html")
	}

	if strings.HasPrefix(path, "static/") {
		// w.Header().Set("Cache-Control", "public, max-age=31536000")
		// c.Request().Response.Header.Set("Cache-Control", "public, max-age=31536000")
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		// w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		// c.Request().Response().Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Logger().Debug(err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))

	return c.Blob(http.StatusOK, contentType, data)
}

func getTask1(c echo.Context) error {

	fmt.Println("in getTask1")
	var res []string

	taskData := <-worm.Out
	fmt.Println("read from worm.Out.")

	for _, v := range taskData {
		res = append(res, v.Json)
	}
	return c.JSON(http.StatusOK, res)
}

func preview1(c echo.Context) error {

	var url = c.QueryParam("url")

	charpters, err := worm.GetchapterURL(url, worm.MODE_PREVIWE)
	if err != nil {
		fmt.Println("erro in get charpters", err)
	}

	fmt.Println("输出章节：", charpters)

	var charpter10 = ""
	for i := 0; i < 10; i++ {
		charpter10 += charpters[i].ChapterName
		charpter10 += "   "
	}

	res, err := json.Marshal(charpter10)
	if err != nil {
		fmt.Println("erro in write response", err)
	}

	return c.JSON(http.StatusOK, res)
}

func config1(c echo.Context) error {

	fmt.Println("in config")

	configs := worm.GetConfigs()

	c.JSON(http.StatusOK, configs)

	return nil
}

func getnovel(c echo.Context) error {

	fmt.Printf("进入getnovel")

	url := c.QueryParam("url")
	fmt.Println("url is ", url)
	// var url = r.URL.Query()["url"][0]

	worm.GetNovel(url)

	return c.String(http.StatusOK, "success")

}
