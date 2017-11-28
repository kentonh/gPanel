// Package user is a child of package api to handle api calls concerning users
package user

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ennovar/gPanel/pkg/networking"
)

// Logout function is accessed by an API call from the webhost root
// by accessing /user_logout and sending it an empty POST request. This function will
// delete the user-auth cookie session store
func Logout(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var store networking.Store
	if strings.Contains(dir, "bundles/") {
		store = networking.GetStore(networking.ACCOUNT_USER_AUTH)
	} else {
		store = networking.GetStore(networking.SERVER_USER_AUTH)
	}

	err := store.Delete(res, req)

	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, http.StatusText(500), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
