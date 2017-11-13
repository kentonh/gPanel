# Server API Documentation

```go
/*
Relative API Path:
  api/server/status
Request:
  N/A
Response(200, 405):
  "0" - OFF,
  "1" - ON,
  "2" - "MAINTENANCE",
  "3" - "RESTARTING"
*/
func Status(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/start
Request:
  N/A
Response(204, 405, 409):
  N/A
*/
func Start(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/shutdown
Request:
  {
    "graceful": boolean
  }
Response(204, 404, 405):
  N/A
*/
func Shutdown(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/restart
Request:
  {
    "graceful": boolean
  }
Response(204, 404, 405):
  N/A
*/
func Restart(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/maintenance
Request:
  N/A
Response(204, 405):
  N/A
*/
func Maintenance(res http.ResponseWriter, req *http.Request) bool {}
```
