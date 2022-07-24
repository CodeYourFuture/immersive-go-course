# HTTP & Authentication

In this project you're going toearn about long-lived processes, some simple networking and the basics of HTTP.

Timebox: 6 days

Learning objectives:

- Use Go's net/http package to build start a simple server that responds to local requests
- Get to know HTTP GET and response codes
- Get familiar with cURL
- Define URL, header, body and content-type
- Accept parameters in via GET in the query string
- Accept data via a POST request
- Setup authentication via a basic HTTP auth
- Switch to using JWTs
- Accept multiple forms of authentication
- Write tests for the above

## Project

- `go mod init http-auth`
- create empty main package and main function
- `go run .`
- `import "net/http"`
- Basic server:

```go
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})

	http.ListenAndServe(":8080", nil)
}
```

- User `curl` to interact
- Add handlers such that the following URLs and responses work. Use `http.NotFoundHandler()`

```
> curl -i http://localhost:8080/500
HTTP/1.1 500 Internal Server Error
Date: Sat, 25 Jun 2022 11:16:30 GMT
Content-Length: 21
Content-Type: text/plain; charset=utf-8

Internal server error

> curl -i http://localhost:8080/200
HTTP/1.1 200 OK
Date: Sat, 25 Jun 2022 11:17:17 GMT
Content-Length: 3
Content-Type: text/plain; charset=utf-8

200

> curl -i http://localhost:8080/404
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 25 Jun 2022 11:17:29 GMT
Content-Length: 19

404 page not found
```

- Make the index page at `/` returns some HTML to a `GET` request

```
curl -i http://localhost:8080/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:42:30 GMT
Content-Length: 42

<!DOCTYPE html><html><em>Hello, world</em>%
```

- Make the index page accept `POST` requests with some HTML, and return that HTML:

```
curl -i -d "<em>Hi</em>" http://localhost:8080/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:43:20 GMT
Content-Length: 32

<!DOCTYPE html><html><em>Hi</em>
```

- Ensure you've got error handling in the handler function
