package router

import (
	"net/http"
	"strconv"
	"time"
	"github.com/Ennovar/gPanel/pkg/database"
	"net/http/httputil"
)

type Router struct {
	Port int
}

var server http.Server
var domainToPort = make(map[string]int)

func RefreshMap() bool {
	ds, err := database.Open("server/" + database.DB_DOMAINS)
	if err != nil {
		return false
	}
	defer ds.Close()

	var client map[string]database.Struct_Domain

	client, err = ds.ListDomains("*")
	if err != nil {
		return false
	}

	for k, v := range client {
		domainToPort[k] = v.PublicPort
	}

	return true
}

func New() *Router {
	if !RefreshMap() {
		return nil
	}

	r := Router{
		Port: 2080,
	}

	server = http.Server{
		Addr:           "localhost:" + strconv.Itoa(r.Port),
		Handler:        &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				if d, ok := domainToPort[req.Host]; ok {
					req.Header.Set("Host", req.Host)
					req.URL.Scheme = "http"
					req.URL.Host = "127.0.0.1:"+strconv.Itoa(d)
				}
			},
		},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	return &r
}

func (r *Router) Start() {
	go server.ListenAndServe()
}