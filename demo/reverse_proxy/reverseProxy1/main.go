package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var (
	next_addr = "http://127.0.0.1:8012"
	port       = "8011"
)

func reverse(rw http.ResponseWriter, req *http.Request)  {
	//解析下游的地址
	nextUrl,err := url.Parse(next_addr)
	if err != nil {
		log.Print("Parse err:",err)
		return
	}
	fmt.Println("Host:",req.Host)
	fmt.Println("RemoteAddr:",req.RemoteAddr)
	fmt.Println("RequestURI:",req.RequestURI)
	fmt.Println("Method:",req.Method)
	fmt.Println("Header:",req.Header)
	fmt.Println("URL:",*req.URL)
	fmt.Println("Body:",req.Body)
	//复制原请求内容
	req.URL.Scheme = nextUrl.Scheme
	req.URL.Host = nextUrl.Host
	transport := http.DefaultTransport
	//请求下游得到返回值
	rsp, err := transport.RoundTrip(req)
	if err != nil {
		log.Print("RoundTrip err:",err)
		return
	}
	for k, vv := range rsp.Header {
		for _, v := range vv {
			rw.Header().Add(k, v)
		}
	}
	defer rsp.Body.Close()
	//写回上游
	_,_ = bufio.NewReader(rsp.Body).WriteTo(rw)
}

func main() {
	http.HandleFunc("/", reverse)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("ListenAndServe err:",err)
	}
}