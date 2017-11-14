// Package public handles the logic of the public facing website
package public

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ennovar/gPanel/pkg/file"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	DocumentRoot            string
	Port                    int
	GracefulShutdownTimeout time.Duration
	Status                  int
	ClientLogger            *file.Handler
	ServerLogger            *file.Handler
	LoadTimeLogger          *file.Handler
}

var controller Controller
var server http.Server

// New function returns a new PublicWeb type.
func New(root string) *Controller {
	clientLogHandler, _ := file.Open(file.LOG_CLIENT_ERRORS, true, true)
	serverLogHandler, _ := file.Open(file.LOG_CLIENT_ERRORS, true, true)
	loadLogHandler, _ := file.Open(file.LOG_LOADTIME, true, true)

	controller = Controller{
		DocumentRoot: root,
		Port:         3000,
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		ClientLogger:            clientLogHandler,
		ServerLogger:            serverLogHandler,
		LoadTimeLogger:          loadLogHandler,
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
	startTime := time.Now()

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
		path = (con.DocumentRoot + "index.html")
	} else {
		path = (con.DocumentRoot + path)
	}

	f, err := os.Open(path)

	if err != nil {
		con.ClientLogger.Write(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		con.ClientLogger.Write(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
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

	elapsedTime := time.Since(startTime)
	con.LoadTimeLogger.Write(path + " rendered in " + strconv.FormatFloat(elapsedTime.Seconds(), 'f', 6, 64) + " seconds")
}
