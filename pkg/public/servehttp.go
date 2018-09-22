package public

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kentonh/gPanel/pkg/database"
	"github.com/kentonh/gPanel/pkg/routing"
)

func (con *Controller) ServePHP(res http.ResponseWriter, path string) {
	if out, err := exec.Command("php-cgi", path).Output(); err == nil {
		reg := regexp.MustCompile(`Status: (\d{3})`)
		b := bytes.NewReader(out)
		s := bufio.NewScanner(b)
		var status, split int

		for s.Scan() {
			line := s.Text()
			split++

			// Status will be first line, if it exists
			if status == 0 {
				m := reg.FindStringSubmatch(line)
				if len(m) == 0 {
					// Status did not exist, so setting default status
					status = 200
				} else {
					var err error
					status, err = strconv.Atoi(m[1])
					if err != nil {
						status = 200
					}
					continue
				}
			}

			// Blank line between headers and body
			if line == "" {
				split++
				break
			}

			sep := strings.Index(line, ": ")
			res.Header().Add(line[:sep], line[sep+2:])
		}

		if err := s.Err(); err != nil {
			con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
			routing.HttpThrowStatus(http.StatusInternalServerError, res)
			return
		}

		res.WriteHeader(status)
		res.Write([]byte(strings.SplitAfterN(string(out), "\n", split)[split-1]))
		return
	} else {
		con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}

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
	if strings.HasPrefix(req.Host, "www") {
		if len(path) == 0 {
			path = con.Directory + "document_root/" + "index.html"
		} else {
			path = con.Directory + "document_root/" + path
		}
	} else {
		if strings.Count(req.Host, ".") == 2 {
			subdomain := strings.SplitN(req.Host, ".", 2)[0] //Remove sub-domain

			ds, err := database.Open(con.AccountDirectory + database.DB_MAIN)
			if err != nil || ds == nil {
				con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
				routing.HttpThrowStatus(http.StatusInternalServerError, res)
				return
			}

			var sdRoot database.StructSubdomain

			err = ds.Get(database.BUCKET_SUBDOMAINS, []byte(subdomain), &sdRoot)
			if err != nil {
				con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
				routing.HttpThrowStatus(http.StatusInternalServerError, res)
				return
			}

			_ = ds.Close()

			if len(path) == 0 {
				path = con.Directory + "document_root/" + sdRoot.Root + "/index.html"
			} else {
				path = con.Directory + "document_root/" + sdRoot.Root + "/" + path
			}
		} else {
			if len(path) == 0 {
				path = con.Directory + "document_root/" + "index.html"
			} else {
				path = con.Directory + "document_root/" + path
			}
		}
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		if strings.HasSuffix(path, ".php") {
			con.ServePHP(res, path)
		} else {
			con.PublicLogger.Println(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
			routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
			return
		}
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
