package ip

import "net/http"

func Unfilter(res http.ResponseWriter, req *http.Request, dir string) bool {
	return true
}
