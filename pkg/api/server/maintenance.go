// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"net/http"

	"github.com/Ennovar/gPanel/pkg/public"
)

func Maintenance(res http.ResponseWriter, req *http.Request, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	publicServer.Maintenance()
	res.WriteHeader(http.StatusNoContent)

	return true
}
