package main

import (
	"log"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/webhost"
	"github.com/gorilla/context"
)

func main() {
	webhost := webhost.New()

	log.Printf("To Exit: CTRL+C")

	go func() {
		log.Print("Listening (public) on localhost:3000, serving out of the document_roots/public/ directory...")
		_ = webhost.Public.Start()
	}()

	log.Print("Listening (private) on localhost:2082, serving out of the document_roots/webhost/ directory...")
	http.ListenAndServe("localhost:2082", context.ClearHandler(&webhost))
}
