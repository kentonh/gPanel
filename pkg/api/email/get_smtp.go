package email

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
)

func GetSMTP(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "GET" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	ds, err := database.Open(dir + database.DB_SETTINGS)
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
