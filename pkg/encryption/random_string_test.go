// Encryption package has functions inside of it that utilize various encypting and hashing techniques
package encryption

import "testing"

func TestRandomString(t *testing.T) {
	testData := []struct {
		length int
		output string
	}{
		{16, ""},
		{16, ""},
		{16, ""},
		{16, ""},
	}

	for i := 0; i < len(testData); i++ {
		testData[i].output = RandomString(testData[i].length)
	}

	for i := 0; i < len(testData)-1; i++ {
		compare := testData[i].output

		for ii := i + 1; ii < len(testData); ii++ {
			if compare == testData[ii].output {
				t.Errorf("Random string generator generated two strings with the same value. (%s - %s)", compare, testData[ii].output)
			}
		}
	}
}
