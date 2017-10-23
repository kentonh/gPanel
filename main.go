package main

import (
	"log"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/public"
	"github.com/Ennovar/gPanel/pkg/webhost"
)

func main() {
	webhost := webhost.NewPrivateHost()
	public := public.NewPublicWeb()

	log.Printf("To Exit: CTRL+C")

	go func() {
		log.Print("Listening (private) on localhost:2082, serving out of the document_roots/webhost/ directory...")
		http.ListenAndServe("localhost:2082", &webhost)
	}()

	log.Print("Listening (public) on localhost:3000, serving out of the document_roots/public/ directory...")
	http.ListenAndServe("localhost:3000", &public)
}
