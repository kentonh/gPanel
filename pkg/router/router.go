package router

import (
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"log"
	"sync"

	"github.com/kentonh/gPanel/pkg/database"
)

type Router struct {
	InsecurePort int
	SecurePort   int
}

var secureServer http.Server
var insecureServer http.Server

var domainToPort map[string]database.Struct_Domain
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
	domainToPort = make(map[string]database.Struct_Domain)
	for k, v := range client {
		domainToPort[k] = v
	}
	mutex.Unlock()

	return true
}

func New(insecure, secure int) *Router {
	if !RefreshMap() {
		return nil
	}

	r := Router{
		InsecurePort: insecure,
		SecurePort:   secure,
	}

	insecureServer = http.Server{
		Addr: "localhost:" + strconv.Itoa(r.InsecurePort),
		Handler: &httputil.ReverseProxy{
			Director:  proxyDirectorInsecure,
			Transport: customTrip{},
		},
	}

	secureServer = http.Server{
		Addr: "localhost:" + strconv.Itoa(r.SecurePort),
		Handler: &httputil.ReverseProxy{
			Director:  proxyDirectorSecure,
			Transport: customTrip{},
		},
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
	go insecureServer.ListenAndServe()
	go secureServer.ListenAndServe()
}
