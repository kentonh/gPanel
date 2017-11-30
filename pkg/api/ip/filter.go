// Package ip is an API that deals with the blocking and unblocking of IPs on public servers
package ip

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
)

func Filter(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var blockIPRequestData struct {
		IP   string `json:"ip"`
		Type string `json:"type"`
	}

	err := json.NewDecoder(req.Body).Decode(&blockIPRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	if blockIPRequestData.Type != "maintenance" && blockIPRequestData.Type != "general" {
		logger.Println(req.URL.Path + "::" + " filtered IP type is invalid")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var newIP database.Struct_Filtered_IP
	newIP.IP = blockIPRequestData.IP
	newIP.Type = blockIPRequestData.Type

	err = ds.NewFilteredIP(&newIP)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
