package main

import (
	"fmt"
	"log"
	"myGateway/demo/balance"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main()  {
	lb := balance.LoadBanlanceFactory(balance.LbHash)
	if err := lb.Add("http://127.0.0.1:8015/base"); err != nil {
		log.Println(err)
	}
	if err := lb.Add("http://127.0.0.1:8016/base"); err != nil {
		log.Println(err)
	}
	//拷贝源码中的内容
	proxy := NewSingleHostReverseProxy(lb)
	fmt.Println(http.ListenAndServe("127.0.0.1:8014", proxy))
}

func NewSingleHostReverseProxy(lb balance.LoadBalance) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		nextAddr, err := lb.Next(req.RemoteAddr)
		if err != nil {
			log.Fatal("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

