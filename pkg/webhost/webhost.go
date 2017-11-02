// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/networking"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type PrivateHost struct {
	Directory string
}

// NewPrivateHost returns a new PrivateHost type.
func NewPrivateHost() PrivateHost {
	return PrivateHost{
		Directory: "document_roots/webhost/",
	}
}

// reqAuth function checks to see if the given path requires authentication.
func reqAuth(path string) bool {
	path = strings.ToLower(path)

	dismissibleTypes := []string{".css", ".js"}
	for _, t := range dismissibleTypes {
		if strings.HasSuffix(path, t) {
			return false
		}
	}

	dismissibleFiles := []string{
		"api_testing.html",
		"index.html",
		"user_auth",
		"user_register",
		"user_logut",
	}
	for _, f := range dismissibleFiles {
		if strings.HasSuffix(path, f) {
			return false
		}
	}

	return true
}

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (priv *PrivateHost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (priv.Directory + "index.html")
	} else {
		path = (priv.Directory + path)
	}

	if reqAuth(path) {
		store := networking.GetStore(networking.COOKIES_USER_AUTH)

		session_value, err := store.Read(w, req, "auth")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session_value == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if auth, ok := session_value.(bool); !ok || !auth {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	isApi, _ := api.HandleAPI(path, w, req)

	if isApi {
		// API methods handle HTTP logic from here
		return
	}

	f, err := os.Open(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusNotFound, w)
		logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
		return
	}

	w.Header().Add("Content-Type", contentType)
	_, err = io.Copy(w, f)

	if err != nil {
		routing.HttpThrowStatus(http.StatusInternalServerError, w)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 500 error.")
		return
	}
}
