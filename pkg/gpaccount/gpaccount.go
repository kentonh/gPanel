// Package gpaccount handles the logic of the gPanel account server
package gpaccount

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ennovar/gPanel/pkg/public"
)

type Controller struct {
	Directory               string
	DocumentRoot            string
	Port                    int
	Public                  *public.Controller
	GracefulShutdownTimeout time.Duration
	Status                  int
	AccountLogger           *log.Logger
	APILogger               *log.Logger
}

var controller Controller
var httpserver http.Server

// New returns a new Controller reference.
func New(dir string, accPort int, pubPort int) *Controller {
	f, err := os.OpenFile(dir+"logs/account_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error whilst trying to start server logging instance:", err.Error())
	}

	apiLogger := log.New(f, "API :: ", 3)
	accountLogger := log.New(f, "ACCOUNT :: ", 3)

	controller = Controller{
		Directory:               dir,
		DocumentRoot:            "account/",
		Port:                    accPort,
		Public:                  public.New(dir, pubPort),
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		AccountLogger:           accountLogger,
		APILogger:               apiLogger,
	}

	httpserver = http.Server{
		Addr:           "localhost:" + strconv.Itoa(controller.Port),
		Handler:        &controller,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
	}

	return &controller
}
