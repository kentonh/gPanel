# Server API Documentation

```go
/*
Relative API Path:
  api/server/status
Request:
  N/A
Response():
  N/A
*/
func Status(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/start
Request:
  N/A
Response():
  N/A
*/
func Start(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/Shutdown
Request:
  N/A
Response():
  N/A
*/
func Shutdown(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/Restart
Request:
  N/A
Response():
  N/A
*/
func Restart(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/server/maintenance
Request:
  N/A
Response():
  N/A
*/
func Maintenance(res http.ResponseWriter, req *http.Request) bool {}
```
