package main

import (
	"gPanel/private_server/routing"
	"gPanel/public_server/routing"
	"log"
	"net/http"
)

func main() {
	private := privateRouting.NewPrivateHost()
	public := publicRouting.NewPublicWeb()

	log.Printf("To Exit: CTRL+C")

	go func() {
		log.Print("Listening (private) on localhost:2082, serving out of the private/ directory...")
		http.ListenAndServe("localhost:2082", &private)
	}()

	log.Print("Listening (public) on localhost:3000, serving out of the public/ directory...")
	http.ListenAndServe("localhost:3000", &public)
}
