// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/api/log"
	"github.com/Ennovar/gPanel/pkg/api/server"
	"github.com/Ennovar/gPanel/pkg/api/user"
)

func (con *Controller) apiHandler(res http.ResponseWriter, req *http.Request) (bool, bool) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + "index.html")
	} else {
		path = (con.Directory + path)
	}

	splitUrl := strings.SplitN(path, "api", 2)
	suspectApi := strings.ToLower(splitUrl[len(splitUrl)-1])

	if req.ContentLength > 0 {
		var bundleRequestData struct {
			BName string `json:"bundle_name,omitempty"`
		}

		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return false, false
		}

		b1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		err = json.NewDecoder(b1).Decode(&bundleRequestData)
		if err != nil {
			return false, false
		}

		if specific, ok := con.Bundles[bundleRequestData.BName]; ok {
			switch suspectApi {
			case "/server/status":
				return true, server.Status(res, req, con.APILogger, specific.Public)
			case "/server/start":
				return true, server.Start(res, req, con.APILogger, specific.Public)
			case "/server/shutdown":
				return true, server.Shutdown(res, req, con.APILogger, specific.Public)
			case "/server/maintenance":
				return true, server.Maintenance(res, req, con.APILogger, specific.Public)
			case "/server/restart":
				return true, server.Restart(res, req, con.APILogger, specific.Public)
			case "/log/read":
				return true, log.Read(res, req, con.APILogger, specific.Directory)
			case "/log/delete":
				return true, log.Delete(res, req, con.APILogger, specific.Directory)
			default:
				return false, false
			}
		}
	}

	switch suspectApi {
	case "/user/auth":
		return true, user.Auth(res, req, con.APILogger, con.Directory)
	case "/user/register":
		return true, user.Register(res, req, con.APILogger, con.Directory)
	case "/user/logout":
		return true, user.Logout(res, req, con.APILogger, con.Directory)
	case "/bundle/create":
		return true, bundle.Create(res, req, con.APILogger, con.Bundles)
	case "/bundle/list":
		return true, bundle.List(res, req, con.APILogger, con.Bundles)
	case "/log/read":
		return true, log.Read(res, req, con.APILogger, con.Directory)
	case "/log/delete":
		return true, log.Delete(res, req, con.APILogger, con.Directory)
	default:
		return false, false
	}
}
