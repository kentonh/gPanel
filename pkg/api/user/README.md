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
func UserAuthentication(res http.ResponseWriter, req *http.Request) bool

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
func UserRegistration(res http.ResponseWriter, req *http.Request) bool

/*
Relative API Path:
  api/user/logout
Request:
  N/A
Response (204, 405, 500):
  N/A
*/
func UserLogout(res http.ResponseWriter, req *http.Request) bool
```
