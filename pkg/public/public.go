// Package public handles the logic of the public facing website
package public

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Controller struct {
	Directory               string
	Port                    int
	GracefulShutdownTimeout time.Duration
	Status                  int
	PublicLogger            *log.Logger
	LoadTimeLogger          *log.Logger
}

var controller Controller
var server http.Server

// New function returns a new PublicWeb type.
func New(dir string, port int) *Controller {
	f, err := os.OpenFile(dir+"logs/public_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("Error whilst trying to start public logging instance: %v\n", err.Error())
	}

	fh, err := os.OpenFile(dir+"logs/public_load_time.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("Error whilst trying to start public load time logging instance: %v\n", err.Error())
	}

	publicLogger := log.New(f, "PUBLIC :: ", 3)
	loadLogger := log.New(fh, "LOAD :: ", 3)

	controller = Controller{
		Directory: dir,
		Port:      port,
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		PublicLogger:            publicLogger,
		LoadTimeLogger:          loadLogger,
	}

	server = http.Server{
		Addr:           "localhost:" + strconv.Itoa(controller.Port),
		Handler:        &controller,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
	}

	return &controller
}
