// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"net/http"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/api/user"
)

func (con *Controller) apiHandler(res http.ResponseWriter, req *http.Request, curBundle int) (bool, bool) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + "index.html")
	} else {
		path = (con.Directory + path)
	}

	splitUrl := strings.SplitN(path, "api", 2)
	suspectApi := strings.ToLower(splitUrl[len(splitUrl)-1])

	switch suspectApi {
	case "/user/auth":
		return true, user.Auth(res, req, con.Directory)
	case "/user/register":
		return true, user.Register(res, req, con.Directory)
	case "/user/logout":
		return true, user.Logout(res, req, con.Directory)
	case "/bundle/create":
		return true, bundle.Create(res, req, con.Bundles)
	default:
		return false, false
	}
}
