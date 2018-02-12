// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"log"
	"os"

	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

type Controller struct {
	Directory    string
	DocumentRoot string
	Bundles      map[string]*gpaccount.Controller
	ServerLogger *log.Logger
	APILogger    *log.Logger
}

func New() (*Controller, error) {
	var err error = nil

	f, err := os.OpenFile("server/logs/server_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error whilst trying to start server logging instance:%v", err.Error())
	}

	c := Controller{
		Directory:    "server/",
		DocumentRoot: "document_root/",
		Bundles:      nil,
		ServerLogger: log.New(f, "SERVER :: ", 3),
		APILogger:    log.New(f, "API :: ", 3),
	}

	err = c.detectBundles()
	c.setDefaults()

	return &c, err
}
