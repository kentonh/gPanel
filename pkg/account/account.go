// Package account handles the logic of the gPanel account server
package account

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/file"
	"github.com/Ennovar/gPanel/pkg/public"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory               string
	DocumentRoot            string
	Port                    int
	Public                  *public.Controller
	GracefulShutdownTimeout time.Duration
	Status                  int
	ServerLogger            *file.Handler
}

var controller Controller
var server http.Server

// New returns a new Controller reference.
func New(root string) *Controller {
	serverErrorLogger, _ := file.Open(file.LOG_SERVER_ERRORS, true, true)

	controller = Controller{
		Directory:               root,
		DocumentRoot:            "account/",
		Port:                    2082,
		Public:                  public.New(root + "public/"),
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		ServerLogger:            serverErrorLogger,
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

func (con *Controller) Start() error {
	if con.Status == 1 {
		return errors.New("Account server is already on.")
	}

	con.Status = 1
	go server.ListenAndServe()
	return nil
}

func (con *Controller) Stop(graceful bool) error {
	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := server.Shutdown(context)
		if err == nil {
			return nil
		}

		fmt.Printf("Graceful shutdown failed attempting forced: %v\n", err)
	}

	if err := server.Close(); err != nil {
		return err
	}

	con.Status = 0
	return nil
}

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + "index.html")
	} else {
		path = (con.Directory + path)
	}

	if reqAuth(path) {
		if !con.checkAuth(res, req) {
			con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnauthorized) + "::" + http.StatusText(http.StatusUnauthorized))
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	isApi, _ := api.HandleAPI(res, req, path, con.Public)

	if isApi {
		// API methods handle HTTP logic from here
		return
	}

	f, err := os.Open(path)

	if err != nil {
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
