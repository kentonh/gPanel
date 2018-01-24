package router

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type customTrip struct{}

func (customTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "fail" {
		if req.Body != nil {
			io.Copy(ioutil.Discard, req.Body)
			defer req.Body.Close()
		}

		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     http.StatusText(http.StatusBadRequest),
			Body:       ioutil.NopCloser(&bytes.Reader{}),
		}, nil
	}

	return http.DefaultTransport.RoundTrip(req)
}

func proxyDirector(req *http.Request) {
	host := req.Host
	if strings.Count(host, ".") == 2 {
		host = strings.SplitN(host, ".", 2)[1] //Remove sub-domain
	}

	req.Header.Set("Host", req.Host)
	req.URL.Scheme = "http"

	mutex.Lock()
	if d, ok := domainToPort[host]; ok {
		mutex.Unlock()
		req.URL.Host = "127.0.0.1:" + strconv.Itoa(d)
	} else {
		mutex.Unlock()
		req.URL.Host = "fail"
	}
}
