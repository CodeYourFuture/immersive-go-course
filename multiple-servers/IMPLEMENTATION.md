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

## Apache

Switch static server over to 8083.

`brew install apache2` (or `brew install httpd`?)

`brew services restart httpd`
