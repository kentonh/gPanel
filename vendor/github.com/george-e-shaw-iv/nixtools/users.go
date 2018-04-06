package nixtools

import (
	"os/exec"
	"bytes"
	"errors"
	"strconv"
	"io"
	"os"
	"strings"
)

// Function GetUser takes either an ID or Username
// as an identifier and returns a type User associated
// with the given identifier. createIfNotExist parameter
// can only be utilized if the function is given a
// username, not an ID.
func GetUser(identifier interface{}, createIfNotExist bool) (*User, error) {
	var u User
	var err error

	if id, ok := identifier.(int); ok {
		u.ID = id

		u.Name, err = getUsername(id)
		if err == nil {
			return &u, nil
		}

		return nil, err
	} else if name, ok := identifier.(string); ok {
		u.Name = name

		u.ID, err = getUserID(name)
		if err == nil {
			return &u, nil
		}

		if createIfNotExist {
			u.ID, err = createUser(name)
			if err != nil {
				return nil, err
			}

			return &u, nil
		}

		return nil, err
	}

	return nil, errors.New("identifier must be a valid user id (int) or username (string)")
}

// Function CreateDirectory attempts to create a directory along
// with all of its parent directories (if needed) relative to the
// user's home directory. The perm parameter is in the format of linux
// file permissions (e.g. 0664).
func (u *User) CreateDirectory(dirPath string, perm os.FileMode) (error) {
	return os.MkdirAll("/home/"+u.Name+"/"+dirPath, perm)
}

// Function WriteFile writes to a file with a given mode that
// is set by flags (flag parameter) that are defined within the
// os package. The filePath is relative to the current user's
// home directory. The perm parameter is in the format of linux
// file permissions (e.g. 0664).
func (u *User) WriteFile(filePath string, flag int, perm os.FileMode, content []byte) (error) {
	f, err := os.OpenFile("/home/"+u.Name+"/"+filePath, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(content)
	if err != nil {
		return err
	}

	if n != len(content) {
		f.Close()
		_ = os.Remove("/home/"+u.Name+"/"+filePath)

		return errors.New("full length of content was not able to be written, file was attempted to be deleted")
	}

	return nil
}

// Function DeleteFileOrDirectory attempts to delete a named
// file or directory relative to the user's home directory.
func (u *User) DeleteFileOrDirectory(path string) (error) {
	return os.RemoveAll("/home/"+u.Name+"/"+path)
}

// Function Delete will attempt to delete the given user
// attached to the struct. deleteOwnedFiles parameter
// will remove all files owned by the user and
// removeHome parameter will delete the home directory
// of the user. Warning: parameter deleteOwnedFiles will
// cause the function to take awhile to execute, consider
// using a goroutine to prevent your program from stalling.
func (u *User) Delete(removeHome, deleteOwnedFiles bool) (error) {
	var stderr bytes.Buffer
	var args = []string{
		"--quiet",
	}

	if deleteOwnedFiles {
		args = append(args, "--remove-all-files")
	}
	if removeHome {
		args = append(args, "--remove-home")
	}

	args = append(args, u.Name)

	// Delete the user and try to remove all files associated
	cmd := exec.Command("deluser", args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	// Forcefully remove the users home directory as sometimes
	// the --remove-home or --remove-all-files flags don't work
	if deleteOwnedFiles || removeHome {
		cmd = exec.Command("rm", "-rf", "/home/"+u.Name)
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}

	return nil
}

// Function SetPassword sets the users password
// to the given parameter.
func (u *User) SetPassword(password string) (error) {
	return nil
}

// Function ForcePasswordChange forces a user to change
// their password upon next login.
func (u *User) ForcePasswordChange() (error) {
	cmd := exec.Command("chage", "-d", "0", u.Name)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

// Function Lock lock's a users account, rendering them
// unable to login through any means of authentication.
func (u *User) Lock() (error) {
	cmd := exec.Command("chage", "-E", "0", u.Name)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

// Function Unlock will unlock a users account, allowing
// them to authenticate through all previous mediums.
func (u *User) Unlock() (error) {
	cmd := exec.Command("chage", "-E", "-1", u.Name)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

// Function getUserID returns the ID associated with
// a username if it exists or -1 and an error if the
// given username does not exist.
func getUserID(username string) (int, error) {
	cmd := exec.Command("id", "-u", username)

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return -1, errors.New(stderr.String())
	}

	stdoutParsed := strings.Replace(stdout.String(), "\n", "", -1)

	id, err := strconv.Atoi(stdoutParsed)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Function getUsername returns the username associated
// with an ID if it exists.
func getUsername(id int) (string, error) {
	var stderr, stdout bytes.Buffer
	reader, writer := io.Pipe()

	cmdOut := exec.Command("getent", "passwd", strconv.Itoa(id))
	cmdOut.Stdout = writer

	cmdIn := exec.Command("cut", "-d:", "-f1")
	cmdIn.Stdin = reader
	cmdIn.Stdout = &stdout
	cmdIn.Stderr = &stderr

	cmdOut.Start()
	cmdIn.Start()

	cmdOut.Wait()
	writer.Close()

	cmdIn.Wait()
	reader.Close()

	if len(stderr.Bytes()) > 0 {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

// Function userExists returns true if the given
// username exists within the system or false if
// the given username does not exist.
func UserExists(username string) (bool) {
	cmd := exec.Command("id", "-u", username)

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

// Function createUser will attempt to create a
// new user in the system with the given username.
// Upon success will return the new ID of the
// created user.
func createUser(username string) (int, error) {
	var stderr bytes.Buffer

	args := []string{
		"--disabled-password",
		"--gecos",
		"",
		username,
	}

	cmd := exec.Command("adduser", args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return 0, errors.New(stderr.String())
	}

	id, err := getUserID(username)
	if err != nil {
		return 0, err
	}

	return id, nil
}