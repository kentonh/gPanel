// Package public handles the logic of the public facing website
package public

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory               string
	GracefulShutdownTimeout time.Duration
	Status                  int
}

var controller Controller
var server http.Server

// New function returns a new PublicWeb type.
func New() *Controller {
	controller = Controller{
		Directory:               "document_roots/public/",
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
	}

	server = http.Server{
		Addr:           "localhost:3000",
		Handler:        &controller,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
	}

	return &controller
}

// Start function starts listening on the public server
func (con *Controller) Start() error {
	if con.Status == 1 {
		return errors.New("Public server is already on.")
	}

	con.Status = 1
	go server.ListenAndServe()
	return nil
}

// Stop function stops the server gracefully or forceful, depending on the boolean input
func (con *Controller) Stop(graceful bool) error {
	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := server.Shutdown(context)
		if err != nil {
			fmt.Printf("Graceful shutdown failed attempting forced: %v\n", err)

			err = server.Close()
			if err != nil {
				return err
			}
		}
	}

	err := server.Close()
	if err != nil {
		return err
	}

	con.Status = 0
	return nil
}

// Restart function combines both the start and stop function, using different
// status codes, as it is restarting.
func (con *Controller) Restart(graceful bool) error {
	con.Status = 3

	if graceful {
		context, cancel := context.WithTimeout(context.Background(), con.GracefulShutdownTimeout)
		defer cancel()

		err := server.Shutdown(context)
		if err != nil {
			fmt.Printf("Graceful shutdown failed attempting forced: %v\n", err)

			err = server.Close()
			if err != nil {
				return err
			}
		}
	}

	err := server.Close()
	if err != nil {
		return err
	}

	con.Status = 1
	go server.ListenAndServe()
	return nil
}

func (con *Controller) Maintenance() {
	con.Status = 2
}

// ServeHTTP function routes all requests for the public web server. It is used in the main
// function inside of the http.ListenAndServe() function for the public host.
func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch con.Status {
	case 0: // This will actually never show because this function won't run if the server is off
		http.Error(res, "The server is currently down and not serving requests.", http.StatusServiceUnavailable)
		return
	case 1: // Normal
		break
	case 2: // Maintenance mode
		http.Error(res, "The server is currently maintenance mode and not serving requests.", http.StatusServiceUnavailable)
		return
	case 3: // This will actually never show because this function won't run if the server is off
		http.Error(res, "The server is currently restarting.", http.StatusServiceUnavailable)
		return
	}

	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + "index.html")
	} else {
		path = (con.Directory + path)
	}

	f, err := os.Open(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusNotFound, res)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 500 error.")
		return
	}
}
