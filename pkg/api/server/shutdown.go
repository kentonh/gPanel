// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/public"
)

// Shutdown function is called from /api/server/shutdown and will attempt to shutdown, either gracefully
// or not gracefully (contingent on request data).
func Shutdown(res http.ResponseWriter, req *http.Request, logger *log.Logger, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var shutdownRequestData struct {
		Graceful bool `json:"graceful"`
	}

	err := json.NewDecoder(req.Body).Decode(&shutdownRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = publicServer.Stop(shutdownRequestData.Graceful)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
