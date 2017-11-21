// Package logs is a child of package api to handle api calls concerning log files
package log

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/file"
)

// Read function is accessed from api/logs/read and will attempt to read
// a given log based off of the request data.
func Read(res http.ResponseWriter, req *http.Request, dir string) bool {
	if req.Method != "POST" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var readLogRequestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&readLogRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	var log string
	switch readLogRequestData.Name {
	case "client_errors":
		log = dir + "logs/" + file.LOG_CLIENT_ERRORS
	case "load_time":
		log = dir + "logs/" + file.LOG_LOADTIME
	case "server_errors":
		log = file.LOG_SERVER_ERRORS
	default:
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	handle, err := file.Open(log, true)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer handle.Close(false)

	data, err := handle.Read()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)

	return true
}
