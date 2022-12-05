# Multiple servers

Create file server to serve static HTML files. Create an API server that serves JSON from a database. Run the API and file server as two separate servers. Try to load the website & see CORS issue. Put nginx in front of the file server and the API so they are on a single port and hostname. Learn about how to run services in VMs in the cloud. Replicate this local setup in the cloud on a single VM, with all services running on the same host. Route requests to the service.

Timebox: 5 days

Learning objectives:

- Basic microservices ideas, separating concerns of services
- Configure nginx to talk to 2-3 copies of the API server
- Some web security ideas (CORS)
- Reverse proxy configuration, routing on path
- Health checks

In future:

- Running applications in the cloud
- Using a cloud-hosted databases
- Multi-environment configuration

## Project

### Design

In this project, we'll build something with the following architecture.

![Architecture of this solution](./readme-assets/architecture.png)

We'll make the file and API server, and use open source software called [nginx](https://nginx.org/) as the load balancer and router.

You can follow the arrows to visualise the path the _request_ takes: an arrow from one box to another is getting data. The _response_ path is the arrow reversed.

You will find the words "upstream" and "downstream" used too. Unfortunately this can be confusing because it depends if you are thinking about the request or response path. In general, upstream and downstream are thought of in terms of dependencies, or from the view of a response. So, the file and API servers are "upstream" of the load balancer: data _flows down the stream_ from the file server to the load balancer, and then the browser.

Let's follow an example request, to `http://localhost:8080/index.html`:

1. The browser requests `http://localhost:8080/index.html`
1. The load balancer is listening on this port and receives the HTTP request
1. It looks at the path (`/index.html`) of the request and tried to match it against its configuration
1. The request path does not match `/api/*`, which you can read as "slash api slash anything". The `*` is often called a "wildcard".
1. The request _does_ match `/*` — "slash anything" — so it **routes** the request to the file server

This selection of where to send the request is called routing, and the order we perform the checks matters. `/*` matches everything, so if we checked that first, we'd never send any traffic to our API server. We check in order from the _most specific_ path pattern to the least specific path patterns.

Separating out servers like this is one way real systems are built. Why do we do this? Much of it comes down to **scale**: doing a lot of anything puts strain on computer resources...

#### State

Some code is "stateful" while other code can be "stateless". We often separate code that is stateful from code that is stateless and, as much as possible, reduce the number of stateful systems. This is because state introduces the _possibility_ of incorrectness, failure and data loss, particularly working at scale.

**Stateless** means there is no stored knowledge relating to past requests; each request can be served independently without depending to another. A file server is likely to be stateless: it can serve any file without knowing what other files have been served in the past to a particular client.

Stateless systems scale easily and simply — you just run more of them!

**Stateful** servers store and retrieve information, and requests may depend on each other: for example, a server that handles banking information needs to know how much money is in the account before it can let someone take money out!

This state/stateless split is the common reason for separating a file server from a server that communicates with a database.

#### Different workloads

Sometimes we split code into different servers or systems because there are very different demands on the computer hardware. This is called the "workload" that the code places on the hardware:

- A _CPU-bound workload_ means something that is limited by the speed of the CPU. A task that performs many calculations, running a complex algorithm like video encoding or 3D modelling, is likely to be CPU bound.

- _I/O-bound workloads_ is limited by how fast data can be read or written from disk or the network, and place heavy demands on these. A server that loads and processes many small files is likely to be I/O bound.

- _Memory-bound workloads_ place heavy demands on the amount of memory or RAM the computer has. Workloads that have to a lot of data into memory, such a database or cache server, are likely to be memory bound.

Placing dissimilar workloads on the same computer can force us to buy very expensive and specialised hardware, make scaling difficult, and make each independent workload negatively affect the other.

There's a good, short guide to workloads on [scaleyourapp.com](https://scaleyourapp.com/a-super-helpful-guide-to-understanding-workload-its-types-in-cloud/) which also looks at workloads in terms of usage patterns.

### Module & packages

Our file layout for this project will look like this:

```console
api/
    api.go
assets/
    ... website ...
cmd/
    static-server/
        main.go
    api-server/
        main.go
config/
    nginx.conf
static/
    static.go
go.mod
```

This is because we're building _two_ servers in the same module: `api` and `static`. Each has its own code and functionality.

Each will be similar, but slightly different because one is connecting to a database and the other is serving files.

There will also be command line tools for configuring and starting each server, in the `cmd` directory:

- `go run ./cmd/api-server` — start the API server
- `go run ./cmd/static-server` — start the static server

Specifically, the `cmd/` files will import functionality from `api` and `static` respectively, and run them. This modularity will make the code easier to understand (which is _the most important thing_ for code!). If you need a refesher on modularity in Go, the Go website has [a good guide](https://go.dev/doc/code).

In reality, starting each will look like this:

```console
# api server
$ DATABASE_URL='postgres://localhost:5432/go-server-database' go run ./cmd/api-server --port 8081

# static server
$ go run ./cmd/static-server --path assets --port 8082
```

> 💡 See the [prep README.md](../prep/README.md#command-line-examples) for an explanation of this command line example.

### Static server

Our "static" server will serve the files for a really simple website. The website will fetch images from our API server and display them as an image gallery.

If you have time or simply want to, you can build this website yourself! However, to get us started, here is something that will work.

Below are three files:

- `index.html` — the main page of the website
- `style.css` — stylesheet for the image gallery
- `script.js` — JavaScript that fetches the image from the API and adds them to the page

Put each of these files into a directory called `assets`: we'll tell the static server to serve these files later on.

##### `index.html`

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Image gallery</title>
    <link rel="stylesheet" href="style.css" />
    <script src="script.js" defer></script>
  </head>
  <body>
    <div class="wrapper">
      <div class="content" role="main">
        <h1 class="title">Gallery</h1>
        <h2>Sunsets and animals like you've never seen them before.</h2>
        <div class="gallery">Loading images&hellip;</div>
      </div>
    </div>
  </body>
</html>
```

##### `style.css`

```css
:root {
  --color-bg: #565264;
  --color-main: #ffffff;
  --color-primary: #d6cfcb;
  --color-secondary: #ccb7ae;
  --color-tertiary: #706677;
  --wrapper-height: 87vh;
  --image-max-width: 300px;
  --image-margin: 3rem;
  --font-family: "HK Grotesk";
  --font-family-header: "HK Grotesk";
}

/* Basic page style resets */
* {
  box-sizing: border-box;
}
[hidden] {
  display: none !important;
}

img {
  max-width: 100%;
}

/* Import fonts */
@font-face {
  font-family: HK Grotesk;
  src: url("https://cdn.glitch.me/605e2a51-d45f-4d87-a285-9410ad350515%2FHKGrotesk-Regular.otf?v=1603136326027")
    format("opentype");
}
@font-face {
  font-family: HK Grotesk;
  font-weight: bold;
  src: url("https://cdn.glitch.me/605e2a51-d45f-4d87-a285-9410ad350515%2FHKGrotesk-Bold.otf?v=1603136323437")
    format("opentype");
}

body {
  font-family: HK Grotesk;
  background-color: var(--color-bg);
  color: var(--color-main);
}

/* Page structure */
.wrapper {
  min-height: var(--wrapper-height);
  display: grid;
  place-items: normal center;
  margin: 0 1rem;
}
.content {
  max-width: 1032px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: start;
  justify-content: start;
}

h1 {
  color: var(--color-primary);
  font-style: normal;
  font-weight: bold;
  font-size: 100px;
  line-height: 105%;
  margin: 0;
}

h2 {
  color: var(--color-secondary);
}

.gallery-image img {
  border: 1em solid var(--color-tertiary);
}
```

##### `script.js`

```javascript
function fetchImages(development) {
  if (development) {
    return Promise.resolve([
      {
        title: "Sunset",
        alt_text: "Clouds at sunset",
        url: "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
      },
      {
        title: "Mountain",
        alt_text: "A mountain at sunset",
        url: "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
      },
    ]);
  }
  return fetch("http://localhost:8081/images.json").then((_) => _.json());
}

function timeout(t, v) {
  return new Promise((res) => {
    setTimeout(() => res(v), t);
  });
}

const gallery$ = document.querySelector(".gallery");

fetchImages(true).then(
  (images) => {
    gallery$.textContent = images.length ? "" : "No images available.";

    images.forEach((img) => {
      const imgElem$ = document.createElement("img");
      imgElem$.src = img.url;
      imgElem$.alt = img.alt_text;
      const titleElem$ = document.createElement("h3");
      titleElem$.textContent = img.title;
      const wrapperElem$ = document.createElement("div");
      wrapperElem$.classList.add("gallery-image");
      wrapperElem$.appendChild(titleElem$);
      wrapperElem$.appendChild(imgElem$);
      gallery$.appendChild(wrapperElem$);
    });
  },
  () => {
    gallery$.textContent = "Something went wrong.";
  }
);
```

This code isn't meant to be fancy or be the focus of this exercise. Feel free to improve it (but don't get too distracted doing so)!

#### Static server CLI tool

On to some Go!

We need a `main.go` file in `cmd/static-server/` that calls a `Run` function in `static/`. The `Run` function should, for now, just call `log.Println("Hello!")`.

To do this, we need the `main.go` to know where to find the code. Luckily Go brings all this together in an easy way...

First, our `go.mod` file needs to declare a module name. Let's go with `servers`.

```go
module servers
```

Now, we can start a file in `static/` — let's say `static/static.go` — like this:

```
package static

func Run() {
    // ...
}
```

With this in place, other code in your module can import `servers/static` and use `Run`:

```
package main

import (
    "servers/static"
)
```

The rest is up to you: hook this together and make this work:

```console
$ go run ./cmd/static-server
Hello!
```

Next, we need the CLI tool to know where to look for files.

To do that, add support for a command like flag: `--path` which will be where the static files are read from. We can use the [flag](https://pkg.go.dev/flag) package for this.

Make this work:

```console
$ go run ./cmd/static-server --path assets
path: assets
```

We also want this server to run on a specific port. Make this work:

```console
$ go run ./cmd/static-server --path assets --port 8082
path: assets
port: 8082
```

Remember that it should be `static/static.go` that is doing the printing, not `cmd/static-server/main.go`! The configuration should be passed from one to the other.

#### Static server

Now we've got config being passed forward, we can build the server itself. This will be up to you to figure out!

This is not as complicated as it might sound. Have a look at all the functions in Go's `net/http` package: there's some useful stuff in there. And make sure to read the [`Handle` documentation](https://pkg.go.dev/net/http#ServeMux.Handle) to see how the `net/http` does URL path matching.

It's possible to do this all in <20 lines of code.

At the end, you should be able to run the server and visit [http://localhost:8082](http://localhost:8082) to see the image gallery!

```console
$ go run ./cmd/static-server --path assets --port 8082
```

We aren't loading the list of images from an API yet; they're hard coded in the JavaScript. Making the API work is coming next.

### API server

The API server in this project will be very similar to the one we created in the `server-database` project, if you have completed that one.

This one will again be up to you. Here's what we need:

- A CLI tool at `cmd/api-server/main.go` that collects a `DATABASE_URL` environment variable and `--port` flag, and then runs the API server
- A Postgres database setup with an appropriate schema: `images` with `title`, `url` and `alt_text`, plus a unique ID
- An API server that:
  - Connects to the database
  - Accepts `GET` requests to `/images.json` and responds with JSON
  - Accepts `POST` requests to `/images.json`, adds the image to the database, and responds with JSON
  - Handles errors without exposing the internal details
  - Supports an `indent` query parameter

Don't forget to handle errors and close the database connection.

We don't expose our internal errors directly to the user for a few reasons:

1. It may leak private information (e.g. a database connection string, which may even include a password!), which may be a security risk.
1. It probably isn't useful to them to know.
1. It may contain confusing terminology which may be embarrassing or confusing to expose.

At the end of this part of the project, we should have the following working...

A server that you start like this: `DATABASE_URL='postgres://localhost:5432/go-server-database' go run ./cmd/api-server --port 8081`

We can `curl` the server to `GET` images:

```console
> curl 'http://localhost:8081/images.json?indent=2' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Thu, 11 Aug 2022 20:17:32 GMT
Content-Length: 763

[
  {
    "title": "Sunset",
    "alt_text": "Clouds at sunset",
    "url": "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"
  },
  {
    "title": "Mountain",
    "alt_text": "A mountain at sunset",
    "url": "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"
  },
  {
    "title": "Cat",
    "alt_text": "A cool cat",
    "url": "https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"
  }
]
```

We can `curl` the server to `POST` new images:

```console
> curl 'http://localhost:8081/images.json?indent=2' -i --data '{"title": "Cat", "alt_text": "A cool cat", "url": "https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"}'
HTTP/1.1 200 OK
Content-Type: text/json
Date: Thu, 11 Aug 2022 20:17:32 GMT
Content-Length: 240

{
  "title": "Cat",
  "alt_text": "A cool cat",
  "url": "https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"
}
```

Try to practice Go modularity by splitting up your API server code into modules:

- `api.go` for the DB connection & HTTP handlers — this is the file we already wrote
- `images.go` for all code relating to reading or writing images

### Getting images from the API

We've now built two servers: a static server for files and an API server that reads & writes data from a database.

We can run these at the same time, listening on two different ports: 8082 for the static server and 8081 for the API.

(What happens if you try to run them on the same port? Give this a try if you haven't.)

But our frontend is not yet fetching images from the API server. We'll do that next, but not without running into a bit of a problem.

Update the `script.js` file to talk to the API: update `fetchImages(true)` to `fetchImages(false)`. This will cause the script to load from a URL rather than a static list of images.

However! We've hit a problem. The images won't load, and we can see "Something went wrong."

See if you can debug what's happening here and fix it: check the [developer tools in your browser](https://developer.mozilla.org/en-US/docs/Learn/Common_questions/What_are_browser_developer_tools).

The fix will be a modification to the API server, modifying the response headers.

These are the kinds of issues we often run into when developing a server interacting with other systems, such as a web browser. It's our job to understand and consider how those other systems work when developing.

### Load balancing & routing

In the architecture diagram at the start we had the file and API servers separated, with requests from the browser going through a load balancer and router layer.

This is a common pattern that we find in larger systems. At the most basic level, this layer is acting as a "reverse proxy" for our servers: it is accepting requests, forwarding them on to other servers, and returning responses. Routing refers to this layer sending requests to the appropriate destination according to some criteria, while load balancing refers to distributing requests across multiple instances of a server.

[Here's a good guide to these ideas](https://www.nginx.com/resources/glossary/reverse-proxy-vs-load-balancer), including some information on why we choose to use such a layer.

For our load balancer/proxy we're going to [Nginx](https://www.nginx.com/), which is a very widely used and useful tool for this job.

We're going to run Nginx locally, in our computers, alongside the API and static server:

- When it receives a request to `/api/*` — "anything beginning with `/api/`" — it will forward that request to the API server
- All other requests will go to the static server

First, get Nginx installed by following [this guide](https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/). If you're on macOS, you can use [Homebrew](https://brew.sh) and the [`nginx` formula](https://formulae.brew.sh/formula/nginx#default): `brew install nginx`.

Learning how to configure Nginx end-to-end is out of scope for this course, so here's an _incomplete_ configuration file to get you started. Put this in `config/nginx.conf` folder. Copy the [`mime.types`](./readme-assets/mime.types) to `config/mime.types`.

```conf
# Determines whether nginx should become a daemon (run in the background — daemon – or foreground)
# https://nginx.org/en/docs/ngx_core_module.html#daemon
daemon off;

# For development purposes, log to stderr
# https://nginx.org/en/docs/ngx_core_module.html#error_log
error_log stderr info;

# Defines the number of worker processes. Auto tries to optimise this, likely to the number of CPU cores.
# https://nginx.org/en/docs/ngx_core_module.html#worker_processes
worker_processes auto;

# Directives that affect connection processing.
# https://nginx.org/en/docs/ngx_core_module.html#events
events {
    # Sets the maximum number of simultaneous connections that can be opened by a worker process.
    # https://nginx.org/en/docs/ngx_core_module.html#events
    worker_connections 1024;
}

http {
    include mime.types;

    # Defines the default MIME type of a response.
    # https://nginx.org/en/docs/http/ngx_http_core_module.html#default_type
    default_type text/plain;

    # Log to stdout
    # https://nginx.org/en/docs/http/ngx_http_log_module.html#access_log
    access_log /dev/stdout;

    # Specifies log format.
    # https://nginx.org/en/docs/http/ngx_http_log_module.html#log_format
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
    '$status $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"';

    # By default, NGINX handles file transmission itself and copies the file into the buffer before sending it.
    # Enabling the sendfile directive eliminates the step of copying the data into the buffer and enables direct
    # copying data from one file descriptor to another.
    # https://docs.nginx.com/nginx/admin-guide/web-server/serving-static-content/
    sendfile on;

    # Enable compression
    # https://docs.nginx.com/nginx/admin-guide/web-server/compression/
    gzip on;

    # Sets configuration for a virtual server.
    # https://nginx.org/en/docs/http/ngx_http_core_module.html#server
    server {
        # Port to listen on
        listen 8080;

        # Requests to /api/ are forwarded to a local server running on port 8081
        # https://nginx.org/en/docs/http/ngx_http_core_module.html#location
        location /api/ {
            # proxy_pass [FILL THIS IN]
        }

        # Other request forwarded to a local server running on port 8082
        location / {
            # proxy_pass [FILL THIS IN]
        }
    }
}
```

Once installed, we can run nginx like this:

```console
> nginx -c `pwd`/config/nginx.conf
```

The `-c` argument tells nginx to load a particular config file, rather than its default location.

The config above is incomplete: there is work to do on the `proxy_pass` lines. Follow the nginx documentation to get it working so that `curl http://localhost:8080/` is sent to the static server, but `curl http://localhost:8080/api/images.json` is sent to the API.

### Benchmarking

Now let's test all out using [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html) again.

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

What do you notice about the profile above? How does it compare to what you see?

We can see from the examples above that there is a bit variance in the performance of the requests. Look at the median and max values for `Total:` — they are quite different! Some of our requests are taking a long time.

This issue is happening because our server is struggling to keep up with the all the requests we are sending it: maybe it's running out of CPU or memory to handle so many requests!

Let's see what we can do about that.

---

Important! Dealing with this in the way we're about it to is unrealistic: we haven't looked into **why** there is such variance. In the real world we'd definitely do that first. And load testing locally, on your computer, is a bad way to do it: it doesn't simulate the kinds of real requests that your server would receive, and it doesn't adequately capture important details like the real hardware and the network that is getting data to and from your server. Also, the load testing actually _uses_ some real CPU and memory that would otherwise be used by the server.

The key takeaway here is: load test with realistic requests using a computer that is similar to the one you'll use to host the real server, and don't load test a server from the same computer it is running on.

---

One way to deal with this performance issue is to run multiple copies of the same server and have the load balancer distribute requests to them. This is why it's called load balancing: the load (requests) to the server is balanced (distributed) across multiple underlying servers.

Let's balance load across 3 copies of the API server: investigate the [upstream](http://nginx.org/en/docs/http/ngx_http_upstream_module.html) module in nginx. Remember that each copy of the API server needs to run on different port!

To test if it's working:

- Make sure your API server prints something whenever it gets a request: for example, `log.Println(r.Method, r.URL.EscapedPath())`
- Run a small `ab`: `ab -n 10 -c 10 "http://127.0.0.1:8080/api/images.json"`
- Observe the server logs: the requests are distributed between the servers!

One of the reasons running a load balancer like nginx is so useful is that is will stop sending requests to an "upstream" server that starts failing. Try this out: turn off one of the servers and run another small `ab`: `ab -n 10 -c 10 "http://127.0.0.1:8080/api/images.json"`.

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

Note `upstream server temporarily disabled while connecting to upstream` — it is automatically spotting this and disabling the server. All of the requests still succeeded, they were just routed to the two remaining servers.

What happens if you turn of _all_ the API servers?
