// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/api/user"
	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/networking"
	"github.com/Ennovar/gPanel/pkg/routing"
	jwt "github.com/dgrijalva/jwt-go"
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
		"index.html",
		"user_auth",
		"user_register",
		"user_logout",
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

		session_value, err := store.Read(w, req, "user")
		if err != nil {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "1")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session_value == nil {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "2")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		username, ok := session_value.(string)
		if !ok {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "3")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		stored_secret, err := user.UserSecret(username)
		if stored_secret == "" {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "4")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		session_value, err = store.Read(w, req, "token")
		if err != nil {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "5")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session_value == nil {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "6")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenString, ok := session_value.(string)
		if !ok {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "7")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// if len(tokenString) < 7 {
		// 	logging.Console("DEBUG::", logging.NORMAL_LOG, "8")
		// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	return
		// }
		// tokenString = tokenString[7:]

		keyfunc := func(t *jwt.Token) (interface{}, error) {
			return []byte(stored_secret), nil
		}

		p := jwt.Parser{
			ValidMethods: []string{"HS256", "HS384", "HS512"},
		}
		t, err := p.ParseWithClaims(tokenString, &jwt.StandardClaims{}, keyfunc)

		if err != nil {
			logging.Console("DEBUG::", logging.NORMAL_LOG, username)
			logging.Console("DEBUG::", logging.NORMAL_LOG, tokenString)
			logging.Console("DEBUG::", logging.NORMAL_LOG, stored_secret)
			logging.Console("DEBUG::", logging.NORMAL_LOG, err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims := t.Claims.(*jwt.StandardClaims)
		if claims.Subject != username {
			logging.Console("DEBUG::", logging.NORMAL_LOG, "10")
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
