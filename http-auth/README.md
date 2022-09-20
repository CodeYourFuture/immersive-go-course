# Servers & HTTP requests

In this project we are going to learn about long-lived processes, some networking and the fundamentals of HTTP.

Timebox: 6 days

Learning objectives:

- Use Go's net/http package to build start a simple server that responds to local requests
- Get to know HTTP GET and response codes
- Get familiar with cURL
- Define URL, header, body and content-type
- Accept parameters in via GET in the query string
- Accept data via a POST request
- Setup authentication via basic HTTP auth
- Write tests for the above

## Project

### Making an HTTP server

[Create a new go module](https://go.dev/doc/tutorial/create-module) in this `http-auth` directory: `go mod init http-auth`.

Create empty main package `main.go` and main function. Check it's all working by running the app: `go run .`.

The main library you'll be working with is built-in to Go: `net/http`. Import it for use: `import "net/http"`.

Here's a basic server that we will build from:

```go
package main

import "net/http"

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})

	http.ListenAndServe(":8080", nil)
}
```

Use `curl` to interact with it: `curl -i http://localhost:8080/`

curl is a tool for transfering data from or to a server. It's very useful for testing and interacting with servers that you build.

Using `curl -i` will show us how the server responds, including the response "headers" and "body". The headers contain metadata about the response, such as what type of data is being sent back.

```
> curl -i 'http://localhost:8080/'
HTTP/1.1 200 OK
Date: Sat, 25 Jun 2022 11:17:17 GMT
Content-Length: 25
Content-Type: text/plain; charset=utf-8

Hello, world
```

> 💡 See the [prep README.md](../prep/README.md#command-line-examples) for an explanation of this command line example.

A common [protocol](https://en.m.wikipedia.org/wiki/Communication_protocol) for sending data between clients and servers over the internet is HTTP. It's used for websites, for example.

We can read [lots about HTTP here](https://developer.mozilla.org/en-US/docs/Web/HTTP).

HTTP requests are sent from a client to a server. They come in various types such as `GET`, for reading information, and `POST` for sending information back.

HTTP responses — data sent back to a "client" from a "server" as a result of an HTTP request — can use a set of [standard codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status) to indicate the status of the server or something about the request.

### Status codes

We are going to make a server that responds to `GET` requests with some of the common ones when a client makes a request to the appropriate URL: 200, 404, 500.

Update our go code so that each of the following paths works. Notice how the URL matches the code that is returned:

- `/200` -> 200 OK
- `/404` -> 404 Not found
- `/500` -> 500 Internal Server Error

```
> curl -i 'http://localhost:8080/200'
HTTP/1.1 200 OK
Date: Sat, 25 Jun 2022 11:16:17 GMT
Content-Length: 3
Content-Type: text/plain; charset=utf-8

200

> curl -i 'http://localhost:8080/500'
HTTP/1.1 500 Internal Server Error
Date: Sat, 25 Jun 2022 11:16:30 GMT
Content-Length: 21
Content-Type: text/plain; charset=utf-8

Internal server error

> curl -i 'http://localhost:8080/404'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 25 Jun 2022 11:17:29 GMT
Content-Length: 19

404 page not found
```

Use `http.NotFoundHandler()` for the `404` error.

### The Content-Type header

HTTP requests can return more than just plan text. Next, make the index page at `/` returns some HTML in response to a `GET` request. Make sure the `Content-Type` response header is set: `w.Header().Add("Content-Type", "text/html")`

```
> curl -i 'http://localhost:8080/'
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:42:30 GMT
Content-Length: 42

<!DOCTYPE html><html><em>Hello, world</em>
```

Curl is just one client we can use to make HTTP requests. Take a moment to try out two more that we've already used:

1. A web browser - open up http://localhost:8080/ in Chrome.
2. Postman - make a GET request to http://localhost:8080/ and see the output.

All three of these are clients that know how to speak HTTP, but they do different things with the response data because they have different goals.

The goal of the Content-Type header is to tell the client how it may want to render the response to the user. Try changing the Content-Type header back to `text/plain`, and see what Chrome does with the same response body.

### Methods: GET and POST

Now make the index page accept `POST` requests with some HTML, and return that HTML. You'll need to check the request method: `request.Method`.

```
> curl -i -d "<em>Hi</em>" 'http://localhost:8080/'
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:43:20 GMT
Content-Length: 32

<!DOCTYPE html><html><em>Hi</em>
```

Again, take a look at the response in different clients.

### Query parameters

HTTP requests can also supply "query parameters" in the URL: `/blog-posts?after=2022-05-04`. Make the handler at `/` output the query parameters as a list. Having the output spaced over multiple lines is optional, but done here for readability.

Note that when running commands in a terminal, some characters have special meaning by default, and need escaping - `?` is one of those characters. We've been using single-quotes (`'`s) around all of our URLs because it stops the terminal from making these characters behave specially.

```
> curl -i 'http://localhost:8080?foo=bar'
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

Try putting some HTML into the query params or body. We'll see that it is interpreted as HTML:

```
> curl -i 'http://localhost:8080?foo=<strong>bar</strong>'
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

(Make sure to take a look at this one in a browser!)

This isn't good! This kind of thing can lead to security issues. Search for "XSS attack" to find out more. Let's fix it.

"Escape" the string any time we take some input (data in `POST` or query parameters) and output it back. We'll need to investigate `html.EscapeString(v)`:

```
> curl -i 'http://localhost:8080?foo=<strong>bar</strong>'
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
> curl -i -d "<em>Hi</em>" 'http://localhost:8080/'
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 10:08:21 GMT
Content-Length: 46

<!DOCTYPE html>
<html>
&lt;em&gt;Hi&lt;/em&gt;
```

Take a look at this in a browser too.

### Authentication

Next we're going to add a URL that can only be accessed if we know a username and secret password.

Add an endpoint `/authenticated` that requires the use of [HTTP Basic auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication). It should return a `401 Unauthorized` status code with a `WWW-Authenticate` header if basic auth is not present or does not match a username and password of our choice. Once Basic Auth is provided, it should respond successful!

Go's `http` library comes with some Basic Auth support built-in, so be sure to use it to make the following work:

```
> curl -i 'http://localhost:8080/authenticated'
HTTP/1.1 401 Unauthorized
Www-Authenticate: Basic realm="localhost", charset="UTF-8"
Date: Sun, 24 Jul 2022 14:12:35 GMT
Content-Length: 0
```

```
> curl -i 'http://localhost:8080/authenticated' -H 'Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ='
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 14:13:04 GMT
Content-Length: 38

<!DOCTYPE html>
<html>
Hello username!
```

We can generate the `dXNl...` text [using this website](https://opinionatedgeek.com/Codecs/Base64Encoder). This is "base64 encoded" which we can search for to find a bit more about. Enter `username:password` to get `dXNlcm5hbWU6cGFzc3dvcmQ=`.

It's not a good idea to put secrets like passwords into code (and base64 encoding text doesn't hide it, it just stores it in a different format). So remove any hard-coded usernames and passwords for basic auth, and use `os.Getenv(...)` so that this works:

```
> AUTH_USERNAME=admin AUTH_PASSWORD=long-memorable-password go run .
```

For bonus points, use [a library](https://github.com/joho/godotenv) to support dotenv files, and set your AUTH_USERNAME and AUTH_PASSWORD in a `.env` file.

### Handling load

Next we're going to test how many requests your server can support, and add basic [rate limiting](https://www.cloudflare.com/en-gb/learning/bots/what-is-rate-limiting).

[Follow this guide](https://www.datadoghq.com/blog/apachebench/) to install and use ApacheBench, which will test to see how many requests our server can handle.

```
> ab -n 10000 -c 100 'http://localhost:8080/'

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

If a server receives too many requests at once, it can break (e.g. it may cause the system to run out of memory).

The fact that some of our requests took much longer than others, even though they were doing the same work, suggests that our server was getting stressed. We can see that in the "Percentages of the requests served within a certain time" section - half of the requests took less than 7ms, but the slowest took 53ms - more than 7 times slower.

It's better to protect your server from being asked to handle too many requests than to have it fall over! So use the `rate` library to reject excess requests (> X per second) with a `503 Service Unavailable` error on a `/limited` endpoint.

```
> go get -u golang.org/x/time
```

We will need to import the module:

```go
import "golang.org/x/time/rate"
```

Then create a limiter:

```go
limiter := rate.NewLimiter(100, 30)
```

And use it:

```go
if limiter.Allow() {
    // Respond as normal!
} else {
    // Respond with an error
}
```

If it is working, we will see `Non-2xx responses` and `Failed requests` in our ApacheBench output:

```
> ab -n 100 -c 100 'http://localhost:8080/limited'
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

Notice that all of our requests took about the same time this time around, and none were much slower - this shows that our server wasn't getting stressed.

One of the things we find in real life is that failure is inevitable. Computers lose power, servers get overloaded and slow down or stop working all together, networks break, etc. Our job as engineers isn't to _prevent_ failure, it's to try to make our systems behave as well as possible _depite_ failure.

In this exercise, we chose to make some of our requests fail fast, so that all of the requests that we _did_ process, got processed well (none were really slow, and our server didn't get overloaded).

Through this course, we learnt a lot more about ways we can give users a better experience by controlling _when_ and _how_ things fail.
