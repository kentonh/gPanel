package publicLogging

import (
	"gPanel/general/networking"
	"log"
	"os"
	"strings"
)

const (
	CONSOLE_PREFIX string = "PUBLIC::"
	LOG_FOLDER     string = "logs/public/"
)

/*
  A function to log to the console

  @param logType int
    The integer identifier cooresponding to the type of log.
      The following integers are defined:
        - 1: Normal (log.Println)
        - 2: Fatal (log.Fatal)
  @param msg string
    The message to log

    @return bool
      Dependant on whether or not logType was defined with a valid logType
*/
func Console(logType int, msg string) bool {
	rawClientIP, _ := networking.GetClientIP()
	clientIP := strings.TrimSpace(string(rawClientIP))

	msg = CONSOLE_PREFIX + clientIP + "::" + msg

	switch logType {
	case 1:
		log.Println(msg)
	case 2:
		log.Println(msg)
	default:
		return false
	}

	return true
}

/*
  A function to log errors for the public server into files

  @param fileType int
    The integer identifier cooresponding to the file to log to.
      The following integers are defined:
        - 1: error.log
        - 2: warning.log
  @param msg string
    The message to log.

  @return bool
    Dependant on whether or not fileType was defined with a valid fileType OR if file operations fail
*/
func File(fileType int, msg string) bool {
	msg += "\n"
	fileName := LOG_FOLDER

	switch fileType {
	case 1:
		fileName += "error.log"
	case 2:
		fileName += "warning.log"
	default:
		return false
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	} else {
		_, err := f.Write([]byte(msg))
		if err != nil {
			return false
		} else {
			err := f.Close()
			if err != nil {
				return false
			} else {
				return true
			}
		}
	}
}
