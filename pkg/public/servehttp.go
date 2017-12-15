package public

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"bufio"
	"net/textproto"
	"os/exec"
	"strings"

	"regexp"

	"github.com/Ennovar/gPanel/pkg/routing"
)

var reg = regexp.MustCompile("[0-9]+")

// ServeHTTP function routes all requests for the public web server. It is used in the main
// function inside of the http.ListenAndServe() function for the public host.
func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	if con.Filter(req, "block") {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	switch con.Status {
	case 0: // This will actually never show because this function won't run if the server is off
		http.Error(res, "The server is currently down and not serving requests.", http.StatusServiceUnavailable)
		return
	case 1: // Normal
		break
	case 2: // Maintenance mode
		if !con.Filter(req, "maintenance") {
			http.Error(res, "The server is currently maintenance mode and not serving requests.", http.StatusServiceUnavailable)
			return
		}
	case 3: // This will actually never show because this function won't run if the server is off
		http.Error(res, "The server is currently restarting.", http.StatusServiceUnavailable)
		return
	}

	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = con.Directory + "public/" + "index.html"
	} else {
		path = con.Directory + "public/" + path
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		var out []byte

		if strings.HasSuffix(path, ".php") {
			if out, err = exec.Command("php-cgi", path).Output(); err != nil {
				con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
				routing.HttpThrowStatus(http.StatusInternalServerError, res)
				return
			}
			cgiRes := string(out)
			pos := strings.Index(cgiRes, "; charset=UTF-8") + len("; charset=UTF-8")
			out = []byte(strings.Trim(cgiRes[pos:], "\n"))

			headers := strings.Replace(cgiRes[0:pos], "\n", "\r\n", -1)
			reader := bufio.NewReader(strings.NewReader(headers + "\r\n\r\n"))
			tp := textproto.NewReader(reader)

			mimeHeader, err := tp.ReadMIMEHeader()
			if err != nil {
				con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
				routing.HttpThrowStatus(http.StatusInternalServerError, res)
				return
			}

			httpHeader := http.Header(mimeHeader)
			for k, v := range httpHeader {
				if k == "Status" {
					if code, err := strconv.Atoi(string(reg.Find([]byte(v[0])))); err == nil {
						res.WriteHeader(code)
					}
					continue
				}
				res.Header().Add(k, v[0])
			}
		} else {
			con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
			routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
			return
		}

		res.Write(out)
	} else {
		f, err := os.Open(path)

		if err != nil {
			con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
			routing.HttpThrowStatus(http.StatusNotFound, res)
			return
		}

		res.Header().Add("Content-Type", contentType)

		_, err = io.Copy(res, f)
		if err != nil {
			con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
			routing.HttpThrowStatus(http.StatusInternalServerError, res)
			return
		}
	}

	elapsedTime := time.Since(startTime)
	con.LoadTimeLogger.Println(path + " rendered in " + strconv.FormatFloat(elapsedTime.Seconds(), 'f', 6, 64) + " seconds")
}
