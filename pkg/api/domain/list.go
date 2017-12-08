package domain

import (
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"github.com/Ennovar/gPanel/pkg/database"
)

func List(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var listDomainsReqData struct {
		Bundle string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&listDomainsReqData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open("server/" + database.DB_DOMAINS)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	domains, err := ds.ListDomains(listDomainsReqData.Bundle)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if len(domains) > 0 {
		b, err := json.Marshal(domains)
		if err != nil {
			logger.Println(req.URL.Path + "::" + err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return false
		}

		res.WriteHeader(http.StatusOK)
		res.Write(b)
		return true
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}