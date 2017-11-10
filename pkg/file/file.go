// Package file handles various file operations
package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	LOG_CLIENT_ERRORS = "client_errors.log"
	LOG_SERVER_ERRORS = "server_errors.log"
	LOG_LOADTIME      = "load_time.log"
)

var KNOWN_LOGS = [...]string{LOG_CLIENT_ERRORS, LOG_SERVER_ERRORS, LOG_LOADTIME}

type Handler struct {
	fileHandle *os.File
	path       string
	append     bool
}

func Open(file string, append bool, log bool) (*Handler, error) {
	var err error
	var absPath string
	var f *os.File

	// Generate file path
	if log {
		absPath, err = filepath.Abs("logs/" + file)
	} else {
		absPath, err = filepath.Abs(file)
	}

	// Handle file path errors
	if err != nil {
		return nil, err
	}

	// Open file
	if append {
		f, err = os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	} else {
		f, err = os.OpenFile(absPath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	}

	//Handle open file errors
	if err != nil {
		return nil, err
	}

	// Upon success, return the handler
	return &Handler{
		fileHandle: f,
		path:       absPath,
	}, nil
}

func (h *Handler) checkExistence(createIfNotExist bool) (bool, error) {
	if _, err := os.Stat(h.path); os.IsNotExist(err) {
		if createIfNotExist {
			if h.append {
				h.fileHandle, err = os.OpenFile(h.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			} else {
				h.fileHandle, err = os.OpenFile(h.path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
			}

			if err != nil {
				return false, err
			}
			return false, nil
		} else {
			return false, nil
		}
	}

	return true, nil
}

func (h *Handler) Read() ([]byte, error) {
	_, err := h.checkExistence(true)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(h.fileHandle)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (h *Handler) Write(data string) (int, error) {
	_, err := h.checkExistence(true)
	if err != nil {
		return 0, err
	}

	written, err := h.fileHandle.Write([]byte(data + "\n"))
	if err != nil {
		return 0, err
	}

	return written, err
}

func (h *Handler) Close(delete bool) (error, error) {
	exist, _ := h.checkExistence(false)
	if !exist {
		return nil, nil
	}

	closeErr := h.fileHandle.Close()

	if delete {
		deleteErr := os.Remove(h.path)
		return closeErr, deleteErr
	}

	return closeErr, nil
}
