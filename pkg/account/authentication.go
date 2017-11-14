// Package account handles the logic of the gPanel account server
package account

import (
	"net/http"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/user"
	"github.com/Ennovar/gPanel/pkg/networking"
	jwt "github.com/dgrijalva/jwt-go"
)

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
		"api/user/auth",
		"api/user/register",
		"api/user/logout",
	}
	for _, f := range dismissibleFiles {
		if strings.HasSuffix(path, f) {
			return false
		}
	}

	return true
}

// checkAuth function returns a boolean based on whether or not the current
// caller is authenticated based off of encrypted sessions using JWT values.
func (con *Controller) checkAuth(res http.ResponseWriter, req *http.Request) bool {
	store := networking.GetStore(networking.ACCOUNT_USER_AUTH)

	session_value, err := store.Read(res, req, "user")
	if err != nil || session_value == nil {
		return false
	}

	username, ok := session_value.(string)
	if !ok {
		return false
	}

	stored_secret, err := user.GetSecret(username, con.Directory)
	if stored_secret == "" {
		return false
	}

	session_value, err = store.Read(res, req, "token")
	if err != nil || session_value == nil {
		return false
	}

	tokenString, ok := session_value.(string)
	if !ok {
		return false
	}

	keyfunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(stored_secret), nil
	}

	p := jwt.Parser{
		ValidMethods: []string{"HS256", "HS384", "HS512"},
	}
	t, err := p.ParseWithClaims(tokenString, &jwt.StandardClaims{}, keyfunc)

	if err != nil {
		return false
	}

	claims := t.Claims.(*jwt.StandardClaims)
	if claims.Subject != username {
		return false
	}

	return true
}
