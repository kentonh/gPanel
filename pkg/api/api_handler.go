// Package api handles all API calls
package api

import (
	"net/http"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/log"
	"github.com/Ennovar/gPanel/pkg/api/server"
	"github.com/Ennovar/gPanel/pkg/api/user"
	"github.com/Ennovar/gPanel/pkg/public"
)

// HandleAPI function takes a path and determines if it is an API call, if it is it will
// call the specified API. It returns two booleans, the first being if it is an API call and
// the second is the response of the API call's function.
//
// The path of an API is formatted as such: document_roots/webhost/api/some/thing. The path is split
// by the sequence of characters "api", returning the first half of the URL which is discarded, and the
// second half which will look like /some/thing. The second half is processed to see whether or not it is
// a valid API, and then handled accordingly from there.
func HandleAPI(res http.ResponseWriter, req *http.Request, path string, publicServer *public.Controller) (bool, bool) {
	splitUrl := strings.SplitN(path, "api", 2)
	suspectApi := strings.ToLower(splitUrl[len(splitUrl)-1])

	switch suspectApi {
	case "/user/auth":
		return true, user.Auth(res, req)
	case "/user/register":
		return true, user.Register(res, req)
	case "/user/logout":
		return true, user.Logout(res, req)
	case "/server/status":
		return true, server.Status(res, req, publicServer)
	case "/server/start":
		return true, server.Start(res, req, publicServer)
	case "/server/shutdown":
		return true, server.Shutdown(res, req, publicServer)
	case "/server/restart":
		return true, server.Restart(res, req, publicServer)
	case "/server/maintenance":
		return true, server.Maintenance(res, req, publicServer)
	case "/log/read":
		return true, log.Read(res, req)
	case "/log/delete":
		return true, log.Delete(res, req)
	default:
		return false, false
	}
}
