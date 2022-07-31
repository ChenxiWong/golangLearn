package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func serveReverseProxy(
	target string, res http.ResponseWriter, req *http.Request) {
	Newurl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(Newurl)
	req.URL.Host = Newurl.Host
	req.URL.Scheme = Newurl.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = Newurl.Host
	proxy.ServeHTTP(res, req)
}

func main() {
	http.HandleFunc("/", doRequest)
	err := http.ListenAndServe(":8000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Print("hello world.")

}
func doRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/hello" {
		serveReverseProxy("http://127.0.0.1:8080", w, r)
	} else {
		fmt.Println(r.Method, r.URL)
		w.Write([]byte(r.URL.Path))
	}

}
