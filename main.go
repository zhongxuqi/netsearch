package main

import (
	"net/http"
	"time"
	"DesertEagleSite/handlefunc"
	"net/http/cookiejar"
)

func main() {
	jar, _ := cookiejar.New(nil)
	http.DefaultClient.Jar = jar
	http.DefaultClient.Timeout = 5 * time.Second
	http.HandleFunc("/", handlefunc.HandleMain)

	http.ListenAndServe("0.0.0.0:8089", nil)
}
