// Encryption package has functions inside of it that utilize various encypting and hashing techniques
package encryption

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, plainText []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, plainText)

	if err != nil {
		return false, err
	}

	return true, nil
}
