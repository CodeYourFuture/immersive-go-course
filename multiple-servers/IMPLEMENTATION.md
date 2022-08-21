This file contains notes on the implementation of this project, to serve as a guide for writing the README.md. It's in this file to avoid merge conflicts as the README.md is updated.

# Implementation

## Code organisation

There are a few ways to organise this:

- One module, one package with all servers implemented there
- One module, multiple packages with each server

And an open question about to start each server independently:

- A top-level switch from CLI arguments
- Multiple mains
- `cmd` directory

Read https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091

Initially I'll go with one module, multiple packages, and a `cmd` directory:

```
cmd/
    static-server/
        main.go
    api-server/
        main.go
static/
    static.go
api/
    api.go
```

Running will be `go run ./cmd/static-server`

## Static files

A self-contained website is in `assets`. This is just a simple image gallery that loads images from static configuration. Later on I'll update it to load from a URL.

`cmd/static-server/main.go` accept CLI flag to assets directory, create config, pass to the server

Server listens, reads files when it gets a request. Using `http.ServeFile` — works v well.

## API

Copied over from server-database, with file split up:

- `util.go` for `MarshalWithIndent`
- `images.go` for all images
- `api.go` for the DB connection & HTTP handlers

This has the same setup steps as server-database, so those can be copied over.

## Ports

`cmd` files should allow ports to be configred using `--port`.

## Nginx

Switch static server over to 8082.

`brew install nginx`

```
nginx -c `pwd`/config/nginx.conf
```

## Benchmark

`ab` the API:

```console
> ab -n 5000 -c 25 "http://127.0.0.1:8080/api/images.json"
This is ApacheBench, Version 2.3 <$Revision: 1901567 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)

Completed 500 requests
Completed 1000 requests
Completed 1500 requests
Completed 2000 requests
Completed 2500 requests
Completed 3000 requests
Completed 3500 requests
Completed 4000 requests
Completed 4500 requests
Completed 5000 requests
Finished 5000 requests


Server Software:        nginx/1.23.1
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /api/images.json
Document Length:        4 bytes

Concurrency Level:      25
Time taken for tests:   1.866 seconds
Complete requests:      5000
Failed requests:        0
Total transferred:      885000 bytes
HTML transferred:       20000 bytes
Requests per second:    2680.08 [#/sec] (mean)
Time per request:       9.328 [ms] (mean)
Time per request:       0.373 [ms] (mean, across all concurrent requests)
Transfer rate:          463.26 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   2.0      1      59
Processing:     1    8  24.2      3     322
Waiting:        0    7  20.3      3     322
Total:          1    9  24.3      4     323

Percentage of the requests served within a certain time (ms)
  50%      4
  66%      4
  75%      5
  80%      6
  90%     28
  95%     34
  98%     37
  99%     62
 100%    323 (longest request)
```

And the static server:

```console
> ab -n 5000 -c 25 "http://127.0.0.1:8080/"
This is ApacheBench, Version 2.3 <$Revision: 1901567 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 500 requests
Completed 1000 requests
Completed 1500 requests
Completed 2000 requests
Completed 2500 requests
Completed 3000 requests
Completed 3500 requests
Completed 4000 requests
Completed 4500 requests
Completed 5000 requests
Finished 5000 requests


Server Software:        nginx/1.23.1
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /
Document Length:        607 bytes

Concurrency Level:      25
Time taken for tests:   1.502 seconds
Complete requests:      5000
Failed requests:        0
Total transferred:      4165000 bytes
HTML transferred:       3035000 bytes
Requests per second:    3328.26 [#/sec] (mean)
Time per request:       7.511 [ms] (mean)
Time per request:       0.300 [ms] (mean, across all concurrent requests)
Transfer rate:          2707.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.3      1       3
Processing:     1    7  10.2      3     115
Waiting:        1    6   9.2      3     115
Total:          2    7  10.1      4     116

Percentage of the requests served within a certain time (ms)
  50%      4
  66%      4
  75%      5
  80%      9
  90%     23
  95%     25
  98%     27
  99%     28
 100%    116 (longest request)
```

## Multiple backends

```nginx
http {
    ...

    # Define a group of API servers that nginx can use
    upstream api {
        server localhost:8081;
    }

    ...

        proxy_pass http://api/;

    ...
}
```

Add alternatives:

```nginx
        server localhost:8083;
        server localhost:8084;
```

Run them:

```console
> DATABASE_URL='postgres://localhost:5432/go-server-database' go run ./cmd/api-server --port 8083
> DATABASE_URL='postgres://localhost:5432/go-server-database' go run ./cmd/api-server --port 8084
```

Run a small `ab` (`ab -n 10 -c 10 "http://127.0.0.1:8080/api/images.json"`) and observe the server logs: the requests are distributed between the servers.

Turn off one of the servers and run another small `ab`: `ab -n 10 -c 10 "http://127.0.0.1:8080/api/images.json"`

Look at the `nginx` logs:

```
127.0.0.1 - - [21/Aug/2022:17:07:44 +0100] "GET /api/images.json HTTP/1.0" 200 4 "-" "ApacheBench/2.3"
2022/08/21 17:07:44 [error] 31112#0: *4088 kevent() reported that connect() failed (61: Connection refused) while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://127.0.0.1:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [warn] 31112#0: *4088 upstream server temporarily disabled while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://127.0.0.1:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [error] 31112#0: *4090 kevent() reported that connect() failed (61: Connection refused) while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://[::1]:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [warn] 31112#0: *4090 upstream server temporarily disabled while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://[::1]:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [error] 31113#0: *4092 kevent() reported that connect() failed (61: Connection refused) while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://127.0.0.1:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [warn] 31113#0: *4092 upstream server temporarily disabled while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://127.0.0.1:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [error] 31113#0: *4092 kevent() reported that connect() failed (61: Connection refused) while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://[::1]:8084/images.json", host: "127.0.0.1:8080"
2022/08/21 17:07:44 [warn] 31113#0: *4092 upstream server temporarily disabled while connecting to upstream, client: 127.0.0.1, server: , request: "GET /api/images.json HTTP/1.0", upstream: "http://[::1]:8084/images.json", host: "127.0.0.1:8080"
```

Note `upstream server temporarily disabled while connecting to upstream` — it is automatically spotting this and disabling the server.
