// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/public"
)

// Maintenance function is called from /api/server/maintenance and will place the public server into
// maintenance mode.
func Maintenance(res http.ResponseWriter, req *http.Request, logger *log.Logger, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	publicServer.Maintenance()
	res.WriteHeader(http.StatusNoContent)

	return true
}
