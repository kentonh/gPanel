package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
)

type GPanel struct {
	authentication int
}

func (g *GPanel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if g.authentication != 1 {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	} else {
		path := req.URL.Path[1:]
		log.Println(path)

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
		}
	}
}

func main() {
	g := &GPanel{authentication: 1}

	log.Print("Listening on localhost:2082")
	http.ListenAndServe("localhost:2082", g)
}
