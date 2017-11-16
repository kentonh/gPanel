package webhost

import (
	"testing"
	//"log"
	"net/http/httptest"
	"net/http"
	. "github.com/smartystreets/goconvey/convey"
	"log"
)


// TestWebhost function is used to test different scenarios happened to the webhost package
func TestWebhost(t *testing.T) {

	host := PrivateHost{1, "../../document_roots/webhost/"}

	Convey("Spining off a testing server", t, func() {


		server := httptest.NewServer(&host)
		defer server.Close()



		Convey("Given a HTTP request for /index.html on webhost", func() {
			res, err := http.Get(server.URL + "/index.html")

			if err != nil {
				log.Fatal(err)
			}

			Convey("Then the response status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, 200)
			})
		})



		Convey("Given a HTTP request for /api_testing.html on webhost", func() {
			res, err := http.Get(server.URL + "/api_testing.html")

			if err != nil {
				log.Fatal(err)
			}

			Convey("Then the response status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, 200)
			})
		})


	})
}