# Log API Documentation

```go
/*
Relative API Path:
  api/log/read
Request:
  {
    "name": string
  }
Response(200, 404, 405):
  [log contents]
*/
func Read(res http.ResponseWriter, req *http.Request) bool {}

/*
Relative API Path:
  api/log/delete
Request:
  {
    "name": string
  }
Response(204, 404, 405):
  N/A
*/
func Delete(res http.ResponseWriter, req *http.Request) bool {}
```
