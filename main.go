package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	app "ui"

	worm "worm"
)

func Open(uri string) error {
	cmd := exec.Command("cmd", "/C", "start "+uri)
	return cmd.Run()
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/api/preview", preview)
	r.HandleFunc("/api/config", config)
	r.HandleFunc("/api/getNovel", getnovel)
	r.HandleFunc("/api/tasks", getTasks)
	// r.HandleFunc("/api/check-token", checkToken)
	// r.HandleFunc("/login", handleLogin)

	r.HandleFunc("/", handleStatic)

	s := http.Server{
		Addr:    ":4321",
		Handler: r,
	}

	Open("http://localhost:4321")

	s.ListenAndServe()

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

func getTasks(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(worm.Tasks)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "accessToken,appKey,User-Agent,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.Header().Set("content-type", "application/json")

	w.Write(res)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {

	fmt.Println("进入handleStatic")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := filepath.Clean(r.URL.Path)

	fmt.Println("path: ", path)

	if path == "/" || path == "\\" { // Add other paths that you route on the UI side here
		path = "index.html"
	}

	path = strings.ReplaceAll(path, "\\", "/")

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "\\")

	fmt.Println("path: ", path)

	file, err := uiFS.Open("index.html")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(w, r)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}

func preview(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.RequestURI)

	var urls = r.URL.Query()["url"]

	fmt.Println(urls[0])

	charpters, err := worm.GetchapterURL(urls[0], worm.MODE_PREVIWE)
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

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "accessToken,appKey,User-Agent,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.Header().Set("content-type", "application/json")

	w.Write(res)

}

func config(w http.ResponseWriter, r *http.Request) {
	configs := worm.GetConfigs()

	data, err := json.Marshal(configs)
	if err != nil {
		fmt.Println(err)

	}

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "accessToken,appKey,User-Agent,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.Header().Set("content-type", "application/json")

	w.Write(data)

}

func getnovel(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("进入getnovel")

	var url = r.URL.Query()["url"][0]

	worm.GetNovel(url)

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "accessToken,appKey,User-Agent,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.Header().Set("content-type", "application/json")

	w.Write([]byte("{status:200}"))

}
