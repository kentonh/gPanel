package publicRouting

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/Ennovar/gPanel/public_server/logging"
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

	publicLogging.Console(1, path)
	f, err := os.Open(path)

	if err == nil {
		bufferedReader := bufio.NewReader(f)
		var contentType string

		if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "text/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
			contentType = "image/jpeg"
		} else {
			contentType = "text/plain"
		}

		w.Header().Add("Content Type", contentType)
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))

		publicLogging.Console(2, "Path \""+path+"\" rendered a 404 error: "+http.StatusText(404))
	}
}
