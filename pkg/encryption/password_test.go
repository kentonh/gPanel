// Encryption package has functions inside of it that utilize various encypting and hashing techniques
package encryption

import "testing"

func TestPasswordHashing(t *testing.T) {
	passwords := []struct {
		ok        bool
		plainText string
		hash      string
	}{
		{true, "fd#@$s4$oiahfoij", ""},
		{false, "ashf324yiuf!@#", ""},
		{true, "4892fjsk#@(!!)", ""},
		{false, "fsufh$&*(#(*f))", ""},
	}

	for _, password := range passwords {
		var err error
		if password.ok {
			password.hash, err = HashPassword(password.plainText)
		} else {
			password.hash, err = HashPassword("this will fail")
		}

		if err != nil {
			t.Errorf("Error in password_test using HashPassword func: %s", err.Error())
		}

		err = CheckPassword([]byte(password.hash), []byte(password.plainText))

		if err != nil {
			if password.ok {
				t.Errorf("Error in password_test using CheckPassword func: %s", err.Error())
			}
		}
	}
}
