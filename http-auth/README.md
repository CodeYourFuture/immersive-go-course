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
> curl -i http://localhost:8080/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:42:30 GMT
Content-Length: 42

<!DOCTYPE html><html><em>Hello, world</em>%
```

- Make the index page accept `POST` requests with some HTML, and return that HTML:

```
> curl -i -d "<em>Hi</em>" http://localhost:8080/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:43:20 GMT
Content-Length: 32

<!DOCTYPE html><html><em>Hi</em>
```

- Ensure you've got error handling in the handler function

- Make the handler at `/` output the query parameters as a list. Having the output spaced over multiple lines is optional, but done here for readability.

```
> curl -i http://localhost:8080\?foo=bar
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:55:33 GMT
Content-Length: 96

<!DOCTYPE html>
<html>
<em>Hello, world</em>
<p>Query parameters:
<ul>
<li>foo: [bar]</li>
</ul>
```

- Try putting some HTML into the query params or body to see that it is interpreted as HTML:

```
> curl -i http://localhost:8080\?foo=\<strong\>bar\</strong
\>
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:57:20 GMT
Content-Length: 113

<!DOCTYPE html>
<html>
<em>Hello, world</em>
<p>Query parameters:
<ul>
<li>foo: [<strong>bar</strong>]</li>
</ul>
```

This isn't good! This kind of thing can lead to security issues. Search for "XSS attack" to find out more. Let's fix it.

- "Escape" the string any time you take some input (data in `POST` or query parameters) and output it back:

```
> curl -i http://localhost:8080\?foo=\<strong\>bar\</strong\>
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 10:08:08 GMT
Content-Length: 125

<!DOCTYPE html>
<html>
<em>Hello, world</em>
<p>Query parameters:
<ul>
<li>foo: [&lt;strong&gt;bar&lt;/strong&gt;]</li>
</ul>
```

```
> curl -i -d "<em>Hi</em>" http://localhost:8080/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 10:08:21 GMT
Content-Length: 46

<!DOCTYPE html>
<html>
&lt;em&gt;Hi&lt;/em&gt;
```

- Add an endpoint `/authenticated` that requires the use of HTTP Basic auth. It should return a `401 Unauthorized` status code with a `WWW-Authenticate` header if basic auth is not present or does not match a username and password of your choice. Once Basic Auth is provided, it should respond successful!

```
> curl -i http://localhost:8080/authenticated
HTTP/1.1 401 Unauthorized
Www-Authenticate: Basic realm="localhost", charset="UTF-8"
Date: Sun, 24 Jul 2022 14:12:35 GMT
Content-Length: 0
```

```
> curl -i http://localhost:8080/authenticated -H 'Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ='
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 14:13:04 GMT
Content-Length: 38

<!DOCTYPE html>
<html>
Hello username!
```

You can generate the `dXNl...` text [using this website](https://opinionatedgeek.com/Codecs/Base64Encoder). This is "base64 encoded" which you can search for to find a bit more about. Enter `username:password` to get `dXNlcm5hbWU6cGFzc3dvcmQ=`.

- It's not a good idea to put secrets like passwords into code. So remove any hard-coded usernames and passwords for basic auth, and use `os.Getenv(...)` so that this works:

```
> AUTH_USERNAME=admin AUTH_PASSWORD=long-memorable-password go run .
```
