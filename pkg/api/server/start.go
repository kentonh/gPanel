// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/public"
)

// Start function is called from /api/server/start and turn the public server on.
func Start(res http.ResponseWriter, req *http.Request, logger *log.Logger, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	err := publicServer.Start()
	if err != nil {
		logger.Println(req.URL.Path + "::" + "public server was not able to start" + "::" + err.Error()) // TODO: should log server name?
		http.Error(res, err.Error(), http.StatusConflict)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
