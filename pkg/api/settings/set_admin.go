package settings

import (
	"net/http"
	"log"
	"strconv"
	"github.com/Ennovar/gPanel/pkg/database"
	"encoding/json"
)

func SetAdmin(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var adminSettingsRequestData database.Struct_Admin

	err := json.NewDecoder(req.Body).Decode(&adminSettingsRequestData)
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

	err = ds.Put(database.BUCKET_GENERAL, []byte("admin"), adminSettingsRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}