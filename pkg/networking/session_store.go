// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var key = []byte("GbP=K4#f$khYuZpStK68GyHxGg$4@5K-")

const (
	COOKIES_USER_AUTH = "gpanel-webhost-user-auth"
)

type store struct {
	handle     *sessions.CookieStore
	cookieName string
}

// GetStore function takes a name and either creates/grabs a store with that name.
func GetStore(name string) store {
	sessionStore := store{
		handle:     sessions.NewCookieStore(key),
		cookieName: name,
	}

	return sessionStore
}

// Set function is attached to the store struct and will set a session value inside of the current store.
func (s *store) Set(res http.ResponseWriter, req *http.Request, key string, value interface{}, expire int) error {
	session, err := s.handle.Get(req, s.cookieName)

	if err != nil {
		return err
	}

	session.Values[key] = value
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   expire,
		HttpOnly: true,
	}
	session.Save(req, res)
	return nil
}

// Read function is attached to the store struct and will read a given session value inside of the current store.
func (s *store) Read(res http.ResponseWriter, req *http.Request, key string) (interface{}, error) {
	session, err := s.handle.Get(req, s.cookieName)

	if err != nil {
		return nil, err
	}

	value := session.Values[key]
	return value, nil
}

// Delete function is attached to the store struct and will delete a given session value inside of the current store.
func (s *store) Delete(res http.ResponseWriter, req *http.Request) error {
	session, err := s.handle.Get(req, s.cookieName)

	if err != nil {
		return err
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	session.Save(req, res)
	return nil
}
