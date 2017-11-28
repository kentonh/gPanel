// Package log is a child of package api to handle api calls concerning log files
package log

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/file"
)

// Read function is accessed from api/logs/read and will attempt to read
// a given log based off of the request data.
func Read(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var readLogRequestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&readLogRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	var log string
	switch readLogRequestData.Name {
	case "public_errors":
		log = dir + "logs/" + LOG_PUBLIC_ERRORS
	case "account_errors":
		log = dir + "logs/" + LOG_ACCOUNT_ERRORS
	case "public_load_time":
		log = dir + "logs/" + LOG_PUBLIC_LOAD
	case "server_errors":
		log = LOG_SERVER_ERRORS
	default:
		logger.Println(req.URL.Path + "::unknown log type requested")
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	handle, err := file.Open(log, true)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer handle.Close(false)

	data, err := handle.Read()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)

	return true
}
