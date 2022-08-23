# Multiple servers

Create file server to serve static HTML files. Create an API server that serves JSON from a database. Run the API and file server as two separate servers. Try to load the website & see CORS issue. Put apache in front of the file server and the API so they are on a single port and hostname. Learn about how to run services in VMs in the cloud. Replicate this local setup in the cloud on a single VM, with all services running on the same host. Route requests to the service.

Timebox: 5 days

Learning objectives:

- Basic microservices ideas, separating concerns of services
- Configure apache to talk to 2-3 copies of the API server
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
1. The request path does not match `/api/*`, which you can read ash "slash api slash anything". The `*` is often called a "wildcard".
1. The request _does_ match `/*` — "slash anything" — so it **routes** the request to the file server

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

There will also be command line tools for configuring and starting each server, in the `cmd` directory:

- `go run ./cmd/api-server` — start the API server
- `go run ./cmd/static-server` — start the static server

Each will be similar, but slightly different because one is connecting to a database and the other is serving files.

Specifically, the `cmd/` files will import functionality from `api` and `static` respectively, and run them. This modularity will make the code easier to understand (which is _the most important thing_ for code!). If you need a refesher on modularity in Go, the Go website has [a good guide](https://go.dev/doc/code).

In reality, starting each will look like this:

```console
# api server
$ DATABASE_URL='postgres://localhost:5432/go-server-database' go run ./cmd/api-server --port 8081

# static server
$ go run ./cmd/static-server --path assets --port 8082
```
