package testing

import (
	"testing"
	"net/http"
	"log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPublicDefaultPage(t *testing.T) {
	Convey("Given a HTTP request for /index.html on server side", t, func() {
		res, err := http.Get("http://localhost:3000/index.html")
		if err != nil {
			log.Fatal(err)
		}
		Convey("Then the response status code should be 200", func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})
}

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!

// TODO: HELP WANTED!
