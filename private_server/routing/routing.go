package privateRouting

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
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
	if priv.Auth != 1 {
		w.WriteHeader(401)
		w.Write([]byte("401 - " + http.StatusText(401)))
		privateLog(http.StatusText(401))
	} else {
		path := req.URL.Path[1:]
		path = (priv.Directory + path)

		privateLog(path)
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

			privateLog("Path: " + path)
			privateLog(http.StatusText(404))
		}
	}
}

func privateLog(msg string) {
	log.Println("PRIVATE (PORT 2082):: " + msg)
}
