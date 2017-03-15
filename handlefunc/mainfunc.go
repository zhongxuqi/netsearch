package handlefunc

import (
	"os"
	"time"
	"sync"
	"strconv"
	"net/http"
	"net/url"
	"fmt"
	"html/template"
	"strings"
	"encoding/json"
	. "DesertEagleSite/bean"
)

var iconHandler http.Handler = http.FileServer(http.Dir("html/image"))
var urlFuncMap map[string] func(w http.ResponseWriter, r *http.Request)
func init() {
	urlFuncMap = make(map[string] func(w http.ResponseWriter, r *http.Request))

	initSpider()
	initFile();
	initQiniu();
}

func parseKeyword(r *http.Request, keyname string) (string, bool) {
	paramsMap, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", false
	}
	keyword := paramsMap.Get(keyname)
	if len(keyword) > 0 {
		return keyword, true
	} else {
		return "", false
	}
}

var mux sync.Mutex

func writeLog(r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	t := time.Now()
	year, month, day := t.Date()
	filename := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-"+ strconv.Itoa(day) + ".log"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write([]byte(t.Format(time.UnixDate)+"  "))
	file.Write([]byte(r.Method + "  " + r.RemoteAddr + "  " + r.URL.Path + "  " + r.URL.RawQuery + "\n"))
}

func writeResult(w http.ResponseWriter, r *http.Request, msg string, err error) {
	var response BaseResponse
	if err != nil {
		response.Status = "500"
		response.Message = err.Error()
	} else {
		response.Status = "200"
		response.Message = msg
	}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func HandleMain(w http.ResponseWriter, r *http.Request) {
	writeLog(r)

	// go to file server
	if r.URL.Path == "/" {
		t, err := template.ParseFiles("html/index.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, nil)
		return
	}
	if r.URL.Path == "/favicon.ico" {
		iconHandler.ServeHTTP(w, r);
	}
	if (r.URL.Path[0:5] == "/html") || (r.URL.Path[0:5] == "/data") {
		handleFileServer(w, r);
		return;
	}

	// go json server
	if strings.LastIndex(r.URL.Path, "/") <= 0 {
		return
	}
	url := r.URL.Path[1:]
	if mhandleFunc, ok := urlFuncMap[url]; ok {
		mhandleFunc(w, r)
		return
	}
}
