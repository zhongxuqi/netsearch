package handlefunc

import (
	"net/http"
	"html/template"
	// "fmt"
)

var staticHandler http.Handler = http.FileServer(http.Dir(""))
var urlMap map[string] string
func initFile() {
	urlMap = make(map[string] string)

	urlMap["/html"] = "html/index.html";
	urlMap["/html/"] = "html/index.html";
	urlMap["/html/index"] = "html/index.html";
}

func handleFileServer(w http.ResponseWriter, r *http.Request) {
	if path, ok := urlMap[r.URL.Path]; ok {
		t, _ := template.ParseFiles(path)
		t.Execute(w, nil)
		return;
	}
	staticHandler.ServeHTTP(w, r)
}
