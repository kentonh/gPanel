package domain

import (
	"net/http"
	"log"
	"strconv"
	"github.com/Ennovar/gPanel/pkg/database"
	"encoding/json"
)

func Link(res http.ResponseWriter, req *http.Request, logger *log.Logger, PublicPort int) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var linkDomainReqData struct {
		Domain string `json:"domain"`
		Bundle string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&linkDomainReqData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open("server/"+database.DB_DOMAINS)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var domainDatabaseData database.Struct_Domain
	domainDatabaseData.BundleName = linkDomainReqData.Bundle
	domainDatabaseData.PublicPort = PublicPort

	err = ds.Put(database.BUCKET_DOMAINS, []byte(linkDomainReqData.Domain), domainDatabaseData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
