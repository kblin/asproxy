package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	target := flag.String("url", "https://antismash.secondarymetabolites.org", "URL to proxy to")
	port := flag.Int("port", 8888, "localhost port to bind to")

	flag.Parse()

	proxy, err := NewProxy(*target)
	if err != nil {
		panic(err)
	}

	bindAddr := fmt.Sprintf("localhost:%d", *port)

	http.HandleFunc("/", ProxyHandler(proxy))
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
