package settings

import (
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"github.com/kentonh/gPanel/pkg/emailer"
	"github.com/kentonh/gPanel/pkg/database"
)

func SetSMTP(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "POST" && req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var smtpRequestData database.Struct_SMTP

	err := json.NewDecoder(req.Body).Decode(&smtpRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	// Test Authentication
	_, err = emailer.New(smtpRequestData.Type, emailer.Credentials{
		Username: smtpRequestData.Username,
		Password: smtpRequestData.Password,
		Server:   smtpRequestData.Server,
		Port:     smtpRequestData.Port,
	})

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

	err = ds.Put(database.BUCKET_GENERAL, []byte("smtp"), smtpRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
