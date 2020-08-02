package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main()  {
	nextPath := "http://127.0.0.1:8015/next"
	nextUrl, err := url.Parse(nextPath)
	if err != nil {
		fmt.Println("Parse err:",err)
	}
	proxy := httputil.NewSingleHostReverseProxy(nextUrl)
	fmt.Println(http.ListenAndServe("127.0.0.1:8014", proxy))
}
