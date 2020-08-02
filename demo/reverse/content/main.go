package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main()  {
	nextPath := "http://127.0.0.1:8015/next"
	nextUrl, err := url.Parse(nextPath)
	if err != nil {
		fmt.Println("Parse err:",err)
	}
	//拷贝源码中的内容
	proxy := NewSingleHostReverseProxy(nextUrl)
	fmt.Println(http.ListenAndServe("127.0.0.1:8014", proxy))
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
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
	//添加修改内容方法
	modifyFunc := func(req *http.Response) error {
		if req.StatusCode == 200 {
			oldBody,err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}
			newBody := []byte("modifyFunc " + string(oldBody))
			req.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
			req.ContentLength = int64(len(newBody))
			req.Header.Set("Content-Length", fmt.Sprint(len(newBody)))
		}else{
			return errors.New("StatusCode!=200")
		}
		return nil
	}
	//添加错误处理方法
	errorHandler := func(rw http.ResponseWriter, req *http.Request, err error) {
		_,_ = rw.Write([]byte(err.Error()))
	}
	return &httputil.ReverseProxy{Director: director,ModifyResponse:modifyFunc,ErrorHandler:errorHandler}
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

