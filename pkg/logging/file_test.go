// Package logging handles the various mediums and logic of logging messages and reports.
package logging

import (
	"testing"
)

func TestFile(t *testing.T) {
	testData := []struct {
		File    string
		Message string
	}{
		{"testfile.log", "testtest123"},
		{"testfile.log", "123:432:abc12"},
	}

	for _, data := range testData {
		written, err := File(data.File, data.Message, true)

		if err != nil {
			t.Errorf("File function in package logging threw an error: %v", err.Error())
		}

		wanted := len([]byte(data.Message + "\n"))
		if written != wanted {
			t.Errorf("File function did not write correct number of bytes in file.\nWanted: %d\nGot: %d", wanted, written)
		}
	}
}
