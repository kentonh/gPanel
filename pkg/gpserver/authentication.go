// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kentonh/gPanel/pkg/api/user"
	jwt "github.com/dgrijalva/jwt-go"
)

// reqAuth function checks to see if the given path requires authentication.
func reqAuth(path string) bool {
	path = strings.ToLower(path)

	dismissibleTypes := []string{".css", ".js", ".jpg", ".png", ".ico"}
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
	c, err := req.Cookie("gpanel-server-user-auth")
	if err != nil {
		return false
	}

	data, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return false
	}

	var sessionData struct {
		Username string `json:"Username"`
		Token    string `json:"Token"`
	}

	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		return false
	}

	storedSecret, err := user.GetSecret(sessionData.Username, con.Directory)
	if storedSecret == "" || err != nil {
		return false
	}

	keyfunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(storedSecret), nil
	}

	p := jwt.Parser{
		ValidMethods: []string{"HS256", "HS384", "HS512"},
	}

	t, err := p.ParseWithClaims(sessionData.Token, &jwt.StandardClaims{}, keyfunc)
	if err != nil {
		return false
	}

	claims := t.Claims.(*jwt.StandardClaims)
	if claims.Subject != sessionData.Username {
		return false
	}

	return true
}
