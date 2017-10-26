package routing

import (
	"testing"
)

func TestContentType(t *testing.T) {
	cases := []struct {
		input string
		mime  string
		ok    bool
	}{
		{"/path/to/index.html", "text/html", true},
		{"somepath/file.css", "text/css", true},
		{"test.js", "text/javascript", true},
		{"asdf.png", "image/png", true},
		{"asdfpng", "", false},
		{"/something.jpg", "image/jpeg", true},
		{"/something.jpeg", "image/jpeg", true},
		{"this/should/fail/.png", "", false},
	}

	for _, tc := range cases {
		mime, err := GetContentType(tc.input)
		if err != nil && tc.ok == true {
			t.Fatalf("\nGetContentType(%s): Got an error, expected OK.", tc.input)
		}
		if err == nil && tc.ok == false {
			t.Fatalf("\nGetContentType(%s): Expected error, got: %s", tc.input, tc.mime)
		}
		if mime != tc.mime {
			t.Fatalf("\nGetContentType(%s):\nWant: %s\nGot: %s", tc.input, tc.mime, mime)
		}
	}

}
