// Package api handles all API calls
package api

import (
	"net/http"
	"strings"
)

// HandleAPI function takes a path and determines if it is an API call, if it is it will
// call the specified API. It returns two booleans, the first being if it is an API call and
// the second is the response of the API call's function.
func HandleAPI(path string, res http.ResponseWriter, req *http.Request) (bool, bool) {
	splitUrl := strings.Split(path, "/")
	suspectApi := strings.ToLower(splitUrl[len(splitUrl)-1])

	switch suspectApi {
	case "authentication":
		return true, Authentication(res, req)
	default:
		return false, false
	}
}
