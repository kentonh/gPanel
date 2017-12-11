package settings

import (
	"net/http"
	"log"
	"strconv"
	"github.com/Ennovar/gPanel/pkg/database"
	"encoding/json"
)

func RemoveNameserver(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "DELETE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var nameserverRequestData database.Struct_Nameserver

	err := json.NewDecoder(req.Body).Decode(&nameserverRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open("server/" + database.DB_SETTINGS)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	err = ds.Delete(database.BUCKET_NAMESERVERS, []byte(nameserverRequestData.Nameserver))
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
