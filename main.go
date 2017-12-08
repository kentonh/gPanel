package main

import (
	"log"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/gpserver"
	"github.com/gorilla/context"
	"github.com/Ennovar/gPanel/pkg/router"
)

func main() {
	mains := gpserver.New()
	router := router.New()

	router.Start()

	log.Print("To Exit: CTRL+C")
	log.Print("Domain router is listening on localhost:2080")
	log.Print("Listening (server) on localhost:2082, serving out of the server/document_root/ directory...")
	http.ListenAndServe("localhost:2082", context.ClearHandler(mains))
}
