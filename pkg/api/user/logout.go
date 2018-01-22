// Package user is a child of package api to handle api calls concerning users
package user

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"time"
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

	var sessionName string
	if strings.Contains(dir, "bundles/") {
		sessionName = "gpanel-account-user-auth"
	} else {
		sessionName = "gpanel-server-user-auth"
	}

	http.SetCookie(res, &http.Cookie{
		Name:    sessionName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})

	res.WriteHeader(http.StatusNoContent)
	return true
}
