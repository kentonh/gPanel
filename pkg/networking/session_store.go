// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var key = []byte("GbP=K4#f$khYuZpStK68GyHxGg$4@5K-")

const (
	ACCOUNT_USER_AUTH = "gpanel-account-user-auth"
	SERVER_USER_AUTH  = "gpanel-server-user-auth"
)

type Store struct {
	handle     *sessions.CookieStore
	cookieName string
}

// GetStore function takes a name and either creates/grabs a store with that name.
func GetStore(name string) Store {
	sessionStore := Store{
		handle:     sessions.NewCookieStore(key),
		cookieName: name,
	}

	return sessionStore
}

// Set function is attached to the store struct and will set a session value inside of the current store.
func (s *Store) Set(res http.ResponseWriter, req *http.Request, key string, value interface{}, expire int) error {
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
func (s *Store) Read(res http.ResponseWriter, req *http.Request, key string) (interface{}, error) {
	session, err := s.handle.Get(req, s.cookieName)

	if err != nil {
		return nil, err
	}

	value := session.Values[key]
	return value, nil
}

// Delete function is attached to the store struct and will delete a given session value inside of the current store.
func (s *Store) Delete(res http.ResponseWriter, req *http.Request) error {
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
