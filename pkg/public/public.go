// Package public handles the logic of the public facing website
package public

import (
	"io"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type PublicWeb struct {
	Directory       string
	MaintenanceMode int
}

// NewPublicWeb returns a new PublicWeb type.
func NewPublicWeb() PublicWeb {
	return PublicWeb{
		Directory:       "document_roots/public/",
		MaintenanceMode: 0,
	}
}

// ServeHTTP function routes all requests for the public web server. It is used in the main
// function inside of the http.ListenAndServe() function for the public host.
func (pub *PublicWeb) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if pub.MaintenanceMode == 1 {
		res.WriteHeader(http.StatusServiceUnavailable)
		res.Write([]byte("We are currently in maintenance mode, please check back later."))
	}

	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (pub.Directory + "index.html")
	} else {
		path = (pub.Directory + path)
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
