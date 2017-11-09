// Package logging handles the various mediums and logic of logging messages and reports.
package logging

import (
	"os"
	"path/filepath"
)

const (
	LOG_CLIENT_ERRORS = "client_errors.log"
	LOG_SERVER_ERRORS = "server_errors.log"
	LOG_LOADTIME      = "loadtime.log"
)

// File logs a given message to a log file. The logFile parameter is relative
// to the logs directory. The parameter delete will delete the file after
// writing to it if set to true (primarily used to help write tests).
func File(logFile string, msg string, delete bool) (int, error) {
	msg = msg + "\n"

	absPath, err := filepath.Abs("../../logs/" + logFile)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return 0, err
	}
	writtenBytes, err := f.Write([]byte(msg))

	if err != nil {
		return writtenBytes, err
	}
	err = f.Close()

	if err != nil {
		return writtenBytes, err
	}

	if delete {
		os.Remove(absPath)
	}

	return writtenBytes, nil
}
