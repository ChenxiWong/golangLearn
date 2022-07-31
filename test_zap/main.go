package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"test_zap/log_mgr"
	"time"

	"go.uber.org/zap"
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

var log *zap.Logger
var zap_field zap.Field

func main() {
	log = log_mgr.GetLogger()
	zap_field = zap.Field{
		Interface: map[string]string{
			"traceId": "123ewerfaskkljdasflkndsflkdflkasl",
		},
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		ServiceTow()
		wg.Done()
	}()
	go func() {
		ServiceOne()
		wg.Done()
	}()
	wg.Wait()

}
func ServiceOne() {
	newServeMux := http.ServeMux{}
	newServeMux.HandleFunc("/", doRequest)
	s := &http.Server{
		Addr:           ":8000",
		Handler:        &newServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	// http.HandleFunc("/", doRequest)
	// err := http.ListenAndServe(":8000", nil) //设置监听的端口
	if err != nil {
		msg := fmt.Sprintf("ListenAndServe: %v", err)
		log.Fatal(msg, zap_field)
	}
	fmt.Print("hello world.")
}
func ServiceTow() {
	newServeMux := http.ServeMux{}
	newServeMux.HandleFunc("/", doRequestTwo)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        &newServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	// http.HandleFunc("/", doRequestTwo)
	// err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		msg := fmt.Sprintf("ListenAndServe: %v", err)
		log.Fatal(msg, zap_field)
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
	log.Info("doRequest success", zap_field)

}
func doRequestTwo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	w.Write([]byte(r.URL.Path))
	log.Info("doRequesttTwo success", zap_field)
}
