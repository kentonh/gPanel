package ip

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
)

func List(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var listIPRequestData struct {
		Type string `json:"type"`
	}

	err := json.NewDecoder(req.Body).Decode(&listIPRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	if listIPRequestData.Type != "maintenance" && listIPRequestData.Type != "block" {
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

	list, err := ds.GetFilteredIPs(listIPRequestData.Type)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if len(list) > 0 {
		b, err := json.Marshal(list)
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
