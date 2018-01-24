package router

import (
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"log"
	"sync"

	"github.com/Ennovar/gPanel/pkg/database"
)

type Router struct {
	Port int
}

var server http.Server
var domainToPort map[string]int

var mutex = &sync.Mutex{}

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

	mutex.Lock()
	domainToPort = make(map[string]int)
	for k, v := range client {
		domainToPort[k] = v.PublicPort
	}
	mutex.Unlock()

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
		Addr: "localhost:" + strconv.Itoa(r.Port),
		Handler: &httputil.ReverseProxy{
			Director:  proxyDirector,
			Transport: customTrip{},
		},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	// Start scheduled map refresher
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if !RefreshMap() {
					ticker.Stop()
					log.Fatal("Error refreshing domain/bundle pairing for router")
					return
				}
				log.Print("Successfully refreshed domain/bundle pairing map")
			}
		}
	}()

	return &r
}

func (r *Router) Start() {
	go server.ListenAndServe()
}
