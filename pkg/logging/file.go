// Package logging handles the various mediums and logic of logging messages and reports.
package logging

import "os"

const (
	PRIVATE_LOG_FOLDER = "logs/private/"
	PUBLIC_LOG_FOLDER  = "logs/public/"

	ERROR   = "error.log"
	WARNING = "warning.log"
)

// File logs a given message to a log file
func File(logFolder string, logFile string, msg string) (string, error) {
	msg += "\n"

	f, err := os.OpenFile(logFolder+logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return "", err
	} else {
		writtenBytes, err := f.Write([]byte(msg))

		if err != nil {
			return string(writtenBytes), err
		} else {
			err := f.Close()

			if err != nil {
				return string(writtenBytes), err
			} else {
				return string(writtenBytes), nil
			}

		}

	}

}
