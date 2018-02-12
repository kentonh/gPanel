package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/gpserver"
	"github.com/Ennovar/gPanel/pkg/router"
)

func main() {
	const serverPort = 2082

	const insecurePort = 2080
	const securePort = 2443

	server, err := gpserver.New()
	if err != nil {
		log.Printf("\nA non-fatal error has occurred whilst starting the server: %v\n\n", err.Error())
	}

	router := router.New(insecurePort, securePort)

	if router == nil {
		log.Fatal("Error starting router")
	}
	router.Start()

	log.Print("To Exit: CTRL+C")
	log.Print("Domain router is listening on localhost:" + strconv.Itoa(insecurePort) + " (HTTP) and localhost:" + strconv.Itoa(securePort) + " (HTTPS)")
	log.Print("gPanel Server is listening on localhost:" + strconv.Itoa(serverPort) + ", serving out of the server/document_root/ directory...")
	http.ListenAndServe("localhost:"+strconv.Itoa(serverPort), server)
}
