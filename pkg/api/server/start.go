// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"net/http"

	"github.com/Ennovar/gPanel/pkg/public"
)

func Start(res http.ResponseWriter, req *http.Request, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	err := publicServer.Start()
	if err != nil {
		http.Error(res, err.Error(), http.StatusConflict)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
