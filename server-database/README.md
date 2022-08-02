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
}
```

Import `"encoding/json"` and `Marshal`:

```go
b, err := json.Marshal(data)
```
