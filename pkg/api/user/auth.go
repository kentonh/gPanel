// Package user is a child of package api to handle api calls concerning users
package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
	"github.com/Ennovar/gPanel/pkg/networking"
	jwt "github.com/dgrijalva/jwt-go"
)

// Auth function is accessed by an API call from the webhost root
// by accessing /user_auth and sending it a post request with userRequestData
// struct in JSON format.
func Auth(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var userRequestData struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
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

	var userDatabaseData database.Struct_Users

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)

	if err == database.ErrKeyNotExist {
		logger.Println(req.URL.Path + "::user does not exist.")
		http.Error(res, "User does not exist.", http.StatusUnauthorized)
		return false
	}

	err = encryption.CheckPassword([]byte(userDatabaseData.Pass), []byte(userRequestData.Pass))
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return false
	}

	secret := encryption.RandomString(16)
	userDatabaseData.Secret = secret

	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	now := time.Now().UTC()
	claims := jwt.StandardClaims{
		Subject:   userRequestData.User,
		NotBefore: now.Unix(),
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(24 * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	var store networking.Store
	if strings.Contains(dir, "bundles/") {
		store = networking.GetStore(networking.ACCOUNT_USER_AUTH)
	} else {
		store = networking.GetStore(networking.SERVER_USER_AUTH)
	}

	err = store.Set(res, req, "token", token, (60 * 60 * 24))
	err2 := store.Set(res, req, "user", userRequestData.User, (60 * 60 * 24))
	if err != nil || err2 != nil {
		logger.Println(req.URL.Path + "::" + err.Error() + "::" + err2.Error())
		http.Error(res, http.StatusText(500), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
