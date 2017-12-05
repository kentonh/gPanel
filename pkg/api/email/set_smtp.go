package email

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/emailer"
)

func SetSMTP(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" && req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var smtpRequestData struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Port     int    `json:"port"`
	}

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

	ds, err := database.Open(dir + database.DB_SETTINGS)
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
