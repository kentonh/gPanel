// Package public handles the logic of the public facing website
package public

import (
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
	LoadTimeLogger          *file.Handler
}

var controller Controller
var server http.Server

// New function returns a new PublicWeb type.
func New(root string, port int) *Controller {
	clientLogHandler, _ := file.Open(file.LOG_CLIENT_ERRORS, true, true)
	loadLogHandler, _ := file.Open(file.LOG_LOADTIME, true, true)

	controller = Controller{
		DocumentRoot: root,
		Port:         port,
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		ClientLogger:            clientLogHandler,
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
		fmt.Printf("Server error serving files to client: %v\n", err)
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}

	elapsedTime := time.Since(startTime)
	con.LoadTimeLogger.Write(path + " rendered in " + strconv.FormatFloat(elapsedTime.Seconds(), 'f', 6, 64) + " seconds")
}
