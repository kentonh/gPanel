package publicRouting

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
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

	publicLog(path)
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

		publicLog("Path: " + path)
		publicLog(http.StatusText(404))
	}
}

func publicLog(msg string) {
	log.Println("PUBLIC (PORT 3000):: " + msg)
}
