// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/public"
)

// Status function is called from api/server/status and will return the current status of
// the public server.
func Status(res http.ResponseWriter, req *http.Request, logger *log.Logger, publicServer *public.Controller) bool {
	if req.Method != "GET" && req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(strconv.Itoa(publicServer.Status)))

	return true
}
