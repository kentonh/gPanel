// Package user is a child of package api to handle api calls concerning users
package user

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kentonh/gPanel/pkg/database"
	"github.com/kentonh/gPanel/pkg/encryption"
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

	var sessionName string
	if strings.Contains(dir, "bundles/") {
		sessionName = "gpanel-account-user-auth"
	} else {
		sessionName = "gpanel-server-user-auth"
	}

	var sessionData struct {
		Username string `json:"Username"`
		Token    string `json:"Token"`
	}
	sessionData.Username = userRequestData.User
	sessionData.Token = token

	b, err := json.Marshal(sessionData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	http.SetCookie(res, &http.Cookie{
		Name:  sessionName,
		Value: base64.StdEncoding.EncodeToString(b),
		Path:  "/",
	})

	res.WriteHeader(http.StatusNoContent)
	return true
}
