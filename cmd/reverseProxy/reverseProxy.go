package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	delay   int
	httpUrl string
	port    string
)

// 很简单直接拷贝源代码使用
// go build -o reverseProxy
func main() {
	ParseFlag()
	//设施日志时间戳
	log.SetFlags(log.Ldate | log.Ltime)
	sleepTime := time.Duration(delay) * time.Second

	rpURL, err := url.Parse(httpUrl)
	if err != nil {
		log.Fatal(err)
	}
	reversProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(rpURL)
		},
		ModifyResponse: func(response *http.Response) error {
			log.Printf("begin sleep %ds\n", delay)
			time.Sleep(sleepTime)
			log.Printf("end sleep\n")
			return nil
		},
	}
	frontendProxy := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        &reversProxy,
	}
	defer func(frontendProxy *http.Server) {
		err := frontendProxy.Close()
		if err != nil {
			panic(err)
		}
	}(frontendProxy)
	log.Println("start reverseProxy")
	err = frontendProxy.ListenAndServe()
	if err != nil {
		panic(err)
		return
	}
}

func ParseFlag() {
	flag.IntVar(&delay, "delay", 5, "help message for flagname")
	flag.StringVar(&httpUrl, "httpUrl", "http://127.0.0.1:8545", "help message for user")
	flag.StringVar(&port, "port", "12345", "help message for user")
	flag.Parse()
}
