package settings

import (
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"github.com/kentonh/gPanel/pkg/database"
)

func GetSMTP(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "GET" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	ds, err := database.Open("server/" + database.DB_SETTINGS)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var smtpDbData database.Struct_SMTP

	err = ds.Get(database.BUCKET_GENERAL, []byte("smtp"), &smtpDbData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Remove password
	smtpDbData.Password = ""

	b, err := json.Marshal(smtpDbData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write(b)
	return true
}
