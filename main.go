package main

import (
	"log"
	"net/http"

	"github.com/Ennovar/gPanel/private_server"
	"github.com/Ennovar/gPanel/public_server"
)

func main() {
	private := private_server.NewPrivateHost()
	public := public_server.NewPublicWeb()

	log.Printf("To Exit: CTRL+C")

	go func() {
		log.Print("Listening (private) on localhost:2082, serving out of the private/ directory...")
		http.ListenAndServe("localhost:2082", &private)
	}()

	log.Print("Listening (public) on localhost:3000, serving out of the public/ directory...")
	http.ListenAndServe("localhost:3000", &public)
}
