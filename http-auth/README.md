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

- [Follow this guide](https://www.datadoghq.com/blog/apachebench/) to install and use ApacheBench, which will test to see how many requests your server can handle

```
> ab -n 10000 -c 100 http://localhost:8080/

This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        76 bytes

Concurrency Level:      100
Time taken for tests:   0.779 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1770000 bytes
HTML transferred:       760000 bytes
Requests per second:    12837.71 [#/sec] (mean)
Time per request:       7.790 [ms] (mean)
Time per request:       0.078 [ms] (mean, across all concurrent requests)
Transfer rate:          2219.02 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4   3.2      3      49
Processing:     1    4   3.2      4      49
Waiting:        0    4   3.1      4      49
Total:          5    8   4.5      7      53

Percentage of the requests served within a certain time (ms)
  50%      7
  66%      8
  75%      8
  80%      8
  90%      8
  95%      9
  98%     10
  99%     11
 100%     53 (longest request)
```

- It's better to protect your server from being asked to handle too many requests than to have it fall over! So use the `rate` library to reject excess requests (> X per second) with a `503 Service Unavailable` error on a `/limited` endpoint.

```
> go get -u golang.org/x/time
```

You will need to import the module:

```go
import "golang.org/x/time/rate"
```

Then create a limiter:

```go
lim := rate.NewLimiter(100, 30)
```

And use it:

```go
if lim.Allow() {
    // Respond as normal!
} else {
    // Respond with an error
}
```

If it is working, you will see `Non-2xx responses` and `Failed requests` in your ApacheBench output:

```
> ab -n 100 -c 100 http://localhost:8080/limited
...

Document Path:          /limited
Document Length:        35 bytes

Concurrency Level:      100
Time taken for tests:   0.006 seconds
Complete requests:      100
Failed requests:        70 <----- HERE!
   (Connect: 0, Receive: 0, Length: 70, Exceptions: 0)
Non-2xx responses:      70 <----- HERE!
Total transferred:      17170 bytes
HTML transferred:       2450 bytes
Requests per second:    15544.85 [#/sec] (mean)
Time per request:       6.433 [ms] (mean)
Time per request:       0.064 [ms] (mean, across all concurrent requests)
Transfer rate:          2606.49 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.7      2       4
Processing:     1    1   0.3      1       4
Waiting:        0    1   0.2      1       1
Total:          2    4   0.7      4       5

Percentage of the requests served within a certain time (ms)
  50%      4
  66%      4
  75%      4
  80%      4
  90%      5
  95%      5
  98%      5
  99%      5
 100%      5 (longest request)
```
