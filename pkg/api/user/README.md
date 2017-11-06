# User API Documentation

```go
/*
Relative API Path:
  api/user/auth
Request:
  {
    "user": "test",
    "pass": "test",
  }
Response(204, 400, 401, 405, 500):
  N/A
*/
func Auth(res http.ResponseWriter, req *http.Request) bool

/*
Relative API Path:
  api/user/register
Request:
  {
    "user": "test",
    "pass": "test",
  }
Response (204, 400, 405, 500):
  N/A
*/
func Register(res http.ResponseWriter, req *http.Request) bool

/*
Relative API Path:
  api/user/logout
Request:
  N/A
Response (204, 405, 500):
  N/A
*/
func Logout(res http.ResponseWriter, req *http.Request) bool
```
