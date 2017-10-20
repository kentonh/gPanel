// Package public_server handles the logic of the public facing website
package public_server

import (
	"bufio"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/general/logging"
	"github.com/Ennovar/gPanel/general/routing"
)

type publicWeb struct {
	Directory string
}

func NewPublicWeb() publicWeb {
	pub := publicWeb{}

	pub.Directory = "public/"

	return pub
}

func (pub *publicWeb) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	path = (pub.Directory + path)

	f, err := os.Open(path)

	if err == nil {
		bufferedReader := bufio.NewReader(f)
		contentType, err := routing.GetContentType(path)

		if err == nil {
			w.Header().Add("Content Type", contentType)
			bufferedReader.WriteTo(w)

			logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 200 success.")
		} else {
			routing.HttpThrowStatus(404, w)
			logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
		}

	} else {
		routing.HttpThrowStatus(404, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
	}

}
