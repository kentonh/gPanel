// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/file"
	"github.com/Ennovar/gPanel/pkg/public"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory    string
	Public       *public.Controller
	ServerLogger *file.Handler
}

// New returns a new PrivateHost type.
func New() Controller {
	serverErrorLogger, _ := file.Open(file.LOG_SERVER_ERRORS, true, true)

	return Controller{
		Directory:    "document_roots/webhost/",
		Public:       public.New(),
		ServerLogger: serverErrorLogger,
	}
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
		if !checkAuth(res, req) {
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
