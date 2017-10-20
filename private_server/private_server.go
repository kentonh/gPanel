// Package private_server handles the logic of the private_server
package private_server

import (
	"bufio"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/general/logging"
	"github.com/Ennovar/gPanel/general/routing"
)

type privateHost struct {
	Auth      int
	Directory string
}

func NewPrivateHost() privateHost {
	priv := privateHost{}

	priv.Auth = 1 // Handle Auth
	priv.Directory = "private/"

	return priv
}

func (priv *privateHost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	path = (priv.Directory + path)

	if priv.Auth != 1 {
		routing.HttpThrowStatus(404, w)
		logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 401 error.")
	} else {
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
