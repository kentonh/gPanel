// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"bufio"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type PrivateHost struct {
	Auth      int
	Directory string
}

// NewPrivateHost returns a new PrivateHost type.
func NewPrivateHost() PrivateHost {
	return PrivateHost{
		Auth:      1,
		Directory: "document_roots/webhost/",
	}
}

func (priv *PrivateHost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (priv.Directory + "index.html")
	} else {
		path = (priv.Directory + path)
	}

	if priv.Auth != 1 {
		routing.HttpThrowStatus(404, w)
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
					routing.HttpThrowStatus(404, w)
					logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
				}

			} else {
				routing.HttpThrowStatus(404, w)
				logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
			}

		}

	}

}
