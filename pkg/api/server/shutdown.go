// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/public"
)

func Shutdown(res http.ResponseWriter, req *http.Request, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var shutdownRequestData struct {
		graceful bool `json:"graceful"`
	}

	err := json.NewDecoder(req.Body).Decode(&shutdownRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = publicServer.Stop(shutdownRequestData.graceful)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
