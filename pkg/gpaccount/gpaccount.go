// Package gpaccount handles the logic of the gPanel account server
package gpaccount

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

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
	// ServerLogger            *file.Handler
}

var controller Controller
var httpserver http.Server

// New returns a new Controller reference.
func New(dir string, accPort int, pubPort int) *Controller {
	// serverErrorLogger, _ := file.Open(file.LOG_SERVER_ERRORS, true, true)

	controller = Controller{
		Directory:               dir,
		DocumentRoot:            "account/",
		Port:                    accPort,
		Public:                  public.New(dir+"public/", pubPort),
		GracefulShutdownTimeout: 5 * time.Second,
		Status:                  0,
		// ServerLogger:            serverErrorLogger,
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

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.DocumentRoot + "index.html")
	} else {
		path = (con.DocumentRoot + path)
	}

	if reqAuth(path) {
		if !con.checkAuth(res, req) {
			// con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnauthorized) + "::" + http.StatusText(http.StatusUnauthorized))
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	isApi, _ := con.apiHandler(res, req)

	if isApi {
		// API methods handle HTTP logic from here
		return
	}

	f, err := os.Open(path)

	if err != nil {
		// con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		fmt.Println(err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		// con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		// con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
