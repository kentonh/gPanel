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

func GetStore(name string) store {
	sessionStore := store{
		handle:     sessions.NewCookieStore(key),
		cookieName: name,
	}

	return sessionStore
}

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

func (s *store) Read(res http.ResponseWriter, req *http.Request, key string) (interface{}, error) {
	session, err := s.handle.Get(req, s.cookieName)

	if err != nil {
		return nil, err
	}

	value := session.Values[key]
	return value, nil
}

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
