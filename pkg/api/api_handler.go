// Package api handles all API calls
package api

import (
	"net/http"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/user"
)

// HandleAPI function takes a path and determines if it is an API call, if it is it will
// call the specified API. It returns two booleans, the first being if it is an API call and
// the second is the response of the API call's function.
func HandleAPI(path string, res http.ResponseWriter, req *http.Request) (bool, bool) {
	splitUrl := strings.Split(path, "/")
	suspectApi := strings.ToLower(splitUrl[len(splitUrl)-1])

	switch suspectApi {
	case "api/user/auth":
		return true, user.Auth(res, req)
	case "api/user/register":
		return true, user.Register(res, req)
	case "api/user/logout":
		return true, user.Logout(res, req)
	default:
		return false, false
	}
}
