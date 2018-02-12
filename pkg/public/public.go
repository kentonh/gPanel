// Package public handles the logic of the public facing website
package public

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Ennovar/gPanel/pkg/database"
	"golang.org/x/crypto/acme/autocert"
)

type Controller struct {
	Directory               string
	AccountDirectory        string
	Name                    string
	Port                    int
	GracefulShutdownTimeout time.Duration
	Status                  int
	PublicLogger            *log.Logger
	LoadTimeLogger          *log.Logger
	Server                  http.Server
}

// New function returns a new PublicWeb type.
func New(dir, accountDir, name string, port int) (*Controller, error) {
	ph, lh, err := getLogHandles(dir)
	if err != nil {
		log.Fatalf("Error trying to start logging instances within %v: %v", dir, err.Error())
	}

	controller := Controller{
		Directory:        dir,
		AccountDirectory: accountDir,
		Name:             name,
		Port:             port,
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		PublicLogger:            ph,
		LoadTimeLogger:          lh,
	}

	_ = os.Mkdir("/home/"+controller.Name+"/tls-cache", 0666)

	ds, err := database.Open("server/" + database.DB_DOMAINS)
	if err != nil {
		return nil, err
	}
	defer ds.Close()

	domainMap, err := ds.ListDomains(name)
	if err != nil {
		return nil, err
	}

	var domains []string
	for k := range domainMap {
		domains = append(domains, k)
	}

	certManager := &autocert.Manager{
		Cache:      autocert.DirCache("/home/" + name + "/tls-cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
	}

	controller.Server = http.Server{
		Addr:           "localhost:" + strconv.Itoa(controller.Port),
		Handler:        &controller,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
		TLSConfig:      &tls.Config{GetCertificate: certManager.GetCertificate},
	}

	return &controller, nil
}

// Function getLogHandles returns the handle for the public logger, load logger,
// and, if applicable, an error all in that order.
func getLogHandles(dir string) (*log.Logger, *log.Logger, error) {
	var dirpath, pubpath, loadpath string
	var err error

	if dirpath, err = filepath.Abs(dir + "logs/"); err != nil {
		return nil, nil, err
	}
	if pubpath, err = filepath.Abs(dir + "logs/public_errors.log"); err != nil {
		return nil, nil, err
	}
	if loadpath, err = filepath.Abs(dir + "logs/public_load_time.log"); err != nil {
		return nil, nil, err
	}

	if _, err = os.Stat(dirpath); os.IsNotExist(err) {
		if err := os.Mkdir(dirpath, 0777); err != nil {
			return nil, nil, err
		}
	}

	f, err := os.OpenFile(pubpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, nil, err
	}

	fh, err := os.OpenFile(loadpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, nil, err
	}

	return log.New(f, "PUBLIC :: ", 3), log.New(fh, "LOAD :: ", 3), nil
}
