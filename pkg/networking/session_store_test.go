// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// BUG(george-e-shaw-iv) Says statement coverage for network package is 0.0%, I think this
// has something to do with the fact that I'm trying to test methods appended to a struct.
func TestSessionStore(t *testing.T) {
	storeData := []struct {
		storeName  string
		cookieName string
		key        string
		value      interface{}
	}{
		{"test-store-one", "test-cookie-one", "foo", "bar"},
		{"test-store-two", "test-cookie-two", "baz", true},
		{"test-store-three", "test-cookie-three", "foobar", 32},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		for _, data := range storeData {
			store := GetStore(data.storeName)

			err := store.Set(res, req, data.key, data.value, 60)
			if err != nil {
				t.Errorf("Error in session_store_test: %s", err.Error())
			}

			val, err := store.Read(res, req, data.key)
			if err != nil {
				t.Errorf("Error in session_store_test: %s", err.Error())
			}

			if reflect.TypeOf(data.value) != reflect.TypeOf(val) {
				t.Errorf("Error in session_store_test type checks, expected %s, got %s", reflect.TypeOf(data.value), reflect.TypeOf(val))
			}

			err = store.Delete(res, req)
			if err != nil {
				t.Errorf("Error in session_store_test: %s", err.Error())
			}
		}
	}))
	defer testServer.Close()
}
