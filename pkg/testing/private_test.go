package testing

import (
	"testing"
	"net/http"
	"log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrivateDefaultPage(t *testing.T) {
	Convey("Given a HTTP request for /index.html on client side", t, func() {
		res, err := http.Get("http://localhost:2082/index.html")
		if err != nil {
			log.Fatal(err)
		}
		Convey("Then the response status code should be 200", func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})
}

func TestPrivateEpanelPage(t *testing.T) {
	Convey("Given a HTTP request for /ePanel.html on client side", t, func() {
		res, err := http.Get("http://localhost:2082/ePanel.html")
		if err != nil {
			log.Fatal(err)
		}
		Convey("Then the response status code should be 200", func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})
}

func TestPrivateCheckAuth(t *testing.T) {
	// TODO: HELP WANTED!
}

func TestPrivateCreateUser(t *testing.T) {
	// TODO: HELP WANTED!
}

func TestPrivateGetUser(t *testing.T) {
	// TODO: HELP WANTED!
}

func TestPrivateEditUser(t *testing.T) {
	// TODO: HELP WANTED!
}

func TestPrivateDeleteUser(t *testing.T) {
	// TODO: HELP WANTED!
}

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!
