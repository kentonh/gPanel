// Package user is a child of package api to handle api calls concerning users
package user

import (
	"net/http"

	"github.com/Ennovar/gPanel/pkg/networking"
)

// UserLogout function is accessed by an API call from the webhost root
// by accessing /user_logout and sending it an empty POST request. This function will
// delete the user-auth cookie session store
func UserLogout(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	store := networking.GetStore(networking.COOKIES_USER_AUTH)
	err := store.Delete(res, req)

	if err != nil {
		http.Error(res, http.StatusText(500), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
