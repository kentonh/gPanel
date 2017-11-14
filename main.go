package main

import (
	"log"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/server"
	"github.com/gorilla/context"
)

func main() {
	gpServer := server.New()

	log.Printf("To Exit: CTRL+C")

	log.Print("Listening (server) on localhost:2083, serving out of the server/document_root/ directory...")
	http.ListenAndServe("localhost:2083", context.ClearHandler(gpServer))
}
