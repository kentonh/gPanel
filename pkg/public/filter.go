package public

import (
	"net/http"

	"github.com/kentonh/gPanel/pkg/database"
	"github.com/kentonh/gPanel/pkg/networking"
)

func (con *Controller) Filter(req *http.Request, ftype string) bool {
	ip := networking.GetClientIP(req)

	ds, err := database.Open(con.Directory + database.DB_MAIN)
	if err != nil || ds == nil {
		con.PublicLogger.Println(req.URL.Path + "::" + err.Error())
		return false
	}
	defer ds.Close()

	filtered, err := ds.IsFiltered(ip, ftype)
	if err != nil {
		con.PublicLogger.Println(req.URL.Path + "::" + err.Error())
		return false
	}

	return filtered
}
