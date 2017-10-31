// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"bufio"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/networking"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type PrivateHost struct {
	Directory string
}

// NewPrivateHost returns a new PrivateHost type.
func NewPrivateHost() PrivateHost {
	return PrivateHost{
		Directory: "document_roots/webhost/",
	}
}

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (priv *PrivateHost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (priv.Directory + "index.html")
	} else {
		path = (priv.Directory + path)
	}

	store := networking.GetStore(networking.COOKIES_USER_AUTH)
	val, err := store.Read(w, req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if val != true && !CheckAuth(path) {
		routing.HttpThrowStatus(http.StatusUnauthorized, w)
		logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 401 error.")
	} else {
		isApi, _ := api.HandleAPI(path, w, req)

		if isApi != true {
			f, err := os.Open(path)

			if err == nil {
				bufferedReader := bufio.NewReader(f)
				contentType, err := routing.GetContentType(path)

				if err == nil {
					w.Header().Add("Content Type", contentType)
					bufferedReader.WriteTo(w)

					logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 200 success.")
				} else {
					routing.HttpThrowStatus(http.StatusUnsupportedMediaType, w)
					logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
				}

			} else {
				routing.HttpThrowStatus(http.StatusNotFound, w)
				logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
			}

		}

	}

}
