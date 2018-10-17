package requester

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	heyServer = &http.Server{
		Addr:           ":9876",
		Handler:        &MyServer{},
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 30,
	}
	heyHandlerMap = make(map[string]HandlersFunc)
)

var mimeTypeExt map[string]string = map[string]string{
	".woff2": "application/x-font-woff",
	".woff":  "application/x-font-woff",
}

var zmmService *ZmmService = &ZmmService{}

const (
	STATIC_PREFIX = "static/admin/"
)

type ResultMsg struct {
	Code int
	Msg  string
	Data interface{}
}

type MyServer struct {
}
type HandlersFunc func(http.ResponseWriter, *http.Request)

func (*MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Path
	h := heyHandlerMap[urlStr]
	if h != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
		return
	}
	SendStaticFile(w, r)
}
func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
func SendJson(w http.ResponseWriter, data interface{}) {
	setJsonHeader(w)
	json.NewEncoder(w).Encode(data)
}

func getTimeStr() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func SendFile(w http.ResponseWriter, bytes []byte, filename string) {
	w.Header().Set("Content-Disposition", "attachment;filename=\""+filename+"\"")
	ext := strings.ToLower(path.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = mimeTypeExt[ext]
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(bytes)
}

//读取静态文件
func SendStaticFile(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimLeft(r.URL.Path, "/")
	if urlPath == "" {
		urlPath = "index.html"
	}
	ext := strings.ToLower(path.Ext(urlPath))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = mimeTypeExt[ext]
	}
	w.Header().Set("Content-Type", mimeType)

	file, err := os.Open(STATIC_PREFIX + urlPath)
	if err != nil {
		w.WriteHeader(404)
		io.Copy(w, bytes.NewReader(make([]byte, 1)))
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

//生成参数脚本文件
func CreateShController(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tpl := r.FormValue("tpl")
	AStart := r.FormValue("AStart")
	AEnd := r.FormValue("AEnd")
	BStart := r.FormValue("BStart")
	BEnd := r.FormValue("BEnd")
	option := r.FormValue("option")
	buffer := zmmService.CreateSh(AStart, AEnd, BStart, BEnd, tpl)
	if option == "ok" {
		datas := strings.Split(buffer.String(), "\n")
		SendJson(w, datas)
	}

	if option == "down" {
		SendFile(w, buffer.Bytes(), getTimeStr()+".txt")
	}

}

func initHandler() {
	heyHandlerMap["/CreateSh"] = CreateShController
}

func StartStaticServer() {
	initHandler()
	log.Println("Server starting on port ", heyServer.Addr)
	err := heyServer.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}
