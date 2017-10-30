# Working API Calls

```go
// User Authentication API
/*
JSON Data Required:
  {
    "user": "test",
    "pass": "test",
  }
*/
func UserAuthentication(res http.ResponseWriter, req *http.Request) bool

// User Registration API
/*
JSON Data Required:
  {
    "user": "test",
    "pass": "test",
  }
*/
func UserRegistration(res http.ResponseWriter, req *http.Request) bool
```
