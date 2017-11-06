// Encryption package has functions inside of it that utilize various encypting and hashing techniques
package encryption

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"1234567890!@#$%^&*()"

var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString function takes an integer length value in and returns a
// random string of that size built from the charset constant.
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}
