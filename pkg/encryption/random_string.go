// Encryption package has functions inside of it that utilize various encypting and hashing techniques
package encryption

import (
	"crypto/rand"
	"fmt"
)

func RandomString() (string, error) {
	n := 5
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	s := fmt.Sprintf("%X", b)
	return s, nil
}
