package routing

import (
	"testing"
)

func TestContentType(t *testing.T) {
	cases := []struct {
		input string
		ok    bool
	}{
		{"/path/to/index.html", true},
		{"somepath/file.css", true},
		{"test.js", true},
		{"./relative/asdf.png", true},
		{"/something.jpg", true},
		{"/something.jpeg", true},
		{"doc.pdf", true},
		{"this/should/fail/.png", false},
		{"php/doesnt/have/a/contenttype.php", false},
		{".pdf", false},
		{"asdfpng", false},
	}

	for _, tc := range cases {
		mime, err := GetContentType(tc.input)
		if err != nil && tc.ok == true {
			t.Fatalf("\nGetContentType(%s): Expected OK, got error: %s", tc.input, err)
		}
		if err == nil && tc.ok == false {
			t.Fatalf("\nGetContentType(%s): Expected error, got OK: %s", tc.input, mime)
		}
	}

}
