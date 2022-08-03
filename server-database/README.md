# API, file-server and database

In this project, you'll build another server with simple API that serves data from JSON files, then convert the backend read from a Postgres database, serving data for the API. You'll then turn off the database and learn how to handle errors correctly.

Timebox: 6 days

Learning objectives:

- Build a simple API server that talks JSON
- Understand how a server and a database work together
- Use SQL to read data from a database
- Accept data over a POST request and write it to the database

## Steps

`go mod init server-database`

[intro to Go and JSON](https://go.dev/blog/json)

Create a struct that represents the data:

```go
type Image struct {
	Title   string
	AltText string
	Url     string
}
```

Initialise some data:

```go
data := []Image{
    {"Sunset", "Clouds at sunset", "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"},
    {"Mountain", "A mountain at sunset", "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"},
}
```

Import `"encoding/json"` and `Marshal`:

```go
b, err := json.Marshal(data)
```

Such that this works:

```
> curl http://localhost:8080/images.json -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Wed, 03 Aug 2022 18:06:34 GMT
Content-Length: 487

[{"Title":"Sunset","AltText":"Clouds at sunset","Url":"https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"},{"Title":"Mountain","AltText":"A mountain at sunset","Url":"https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"}]
```

Add a query param `pretty` which uses `MarshalIndent` instead (https://pkg.go.dev/encoding/json#MarshalIndent) such that this works:
