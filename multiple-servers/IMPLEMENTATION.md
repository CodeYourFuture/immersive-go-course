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
