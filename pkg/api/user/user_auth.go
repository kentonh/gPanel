// Package user is a child of package api to handle api calls concerning users
package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
	"github.com/Ennovar/gPanel/pkg/networking"
	jwt "github.com/dgrijalva/jwt-go"
)

// UserAuthentication function is accessed by an API call from the webhost root
// by accessing /user_auth and sending it a post request with userRequestData
// struct in JSON format.
func UserAuthentication(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(database.DBLOC_MAIN)
	if err != nil || ds == nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)

	if err == database.ErrKeyNotExist {
		http.Error(res, "User does not exist.", http.StatusUnauthorized)
		return false
	}

	err = encryption.CheckPassword([]byte(userDatabaseData.Pass), []byte(userRequestData.Pass))
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return false
	}

	secret, err := encryption.RandomString()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	userDatabaseData.Secret = secret

	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
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
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	store := networking.GetStore(networking.COOKIES_USER_AUTH)
	err = store.Set(res, req, "token", token, (60 * 60 * 24))
	err2 := store.Set(res, req, "user", userRequestData.User, (60 * 60 * 24))
	if err != nil || err2 != nil {
		http.Error(res, http.StatusText(500), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
