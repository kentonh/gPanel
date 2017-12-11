package router

import (
	"net/http"
	"strconv"
	"time"
	"github.com/Ennovar/gPanel/pkg/database"
	"net/url"
)

type Router struct {
	Port int
}

var server http.Server

func New() *Router {
	r := Router{
		Port: 2080,
	}

	server = http.Server{
		Addr:           "localhost:" + strconv.Itoa(r.Port),
		Handler:        &r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	return &r
}

func (r *Router) Start() {
	go server.ListenAndServe()
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	domain := req.Host

	ds, err := database.Open("server/" + database.DB_DOMAINS)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ds.Close()

	var client database.Struct_Domain
	err = ds.Get(database.BUCKET_DOMAINS, []byte(domain), &client)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	cURL, err := url.Parse("localhost:"+strconv.Itoa(client.PublicPort)+"/"+req.URL.Path[1:])
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(cURL.String()))
}