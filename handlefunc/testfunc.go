package handlefunc

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type MyData struct {
	Name string `json="name"`
}

func  FooHandler(w http.ResponseWriter, r *http.Request) {
	data := MyData{
		Name: "zhongxuqi",
	}
	b, _ := json.Marshal(data)

	fmt.Fprintf(w, "HTTP/1.1 200 OK\n")
	fmt.Fprintf(w, "Content-Type: text/html \n")
	fmt.Fprintf(w, "Content-Type: %d \n", len(b))
	fmt.Fprintf(w, "\n")
	w.Write(b)
}
