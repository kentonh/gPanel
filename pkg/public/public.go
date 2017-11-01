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
	Directory string
}

// NewPublicWeb returns a new PublicWeb type.
func NewPublicWeb() PublicWeb {
	return PublicWeb{
		Directory: "document_roots/public/",
	}
}

// ServeHTTP function routes all requests for the public web server. It is used in the main
// function inside of the http.ListenAndServe() function for the public host.
func (pub *PublicWeb) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (pub.Directory + "index.html")
	} else {
		path = (pub.Directory + path)
	}

	f, err := os.Open(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusNotFound, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
		return
	}

	w.Header().Add("Content-Type", contentType)
	_, err = io.Copy(w, f)

	if err != nil {
		routing.HttpThrowStatus(http.StatusInternalServerError, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 500 error.")
		return
	}
}
