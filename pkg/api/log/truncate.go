// Package log is a child of package api to handle api calls concerning log files
package log

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Delete function is accessed from api/logs/delete and will attempt to
// delete a given log based off of request data.
func Truncate(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var deleteLogRequestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&deleteLogRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	var log string
	switch deleteLogRequestData.Name {
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

	f, err := os.OpenFile(log, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = f.Close()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
