# Server and database

In this project, we will build another server. This one will have a simple API that serves data in JSON form. We will them convert the backend to read from a Postgres database, serving data for the API. We will then turn off the database and learn how to handle errors correctly.

Timebox: 6 days

Learning objectives:

- Build a simple API server that talks JSON
- Understand how a server and a database work together
- Use SQL to read data from a database
- Accept data over a POST request and write it to the database

## Project

We're going to be using JSON in this project. Before we start, read this [introduction to Go and JSON](https://go.dev/blog/json).

To get started, `cd` into the `server-database` directory, and run `go mod init server-database`.

We've created some `main.go` files in the past: do the same here.

Initially we need to create a struct that represents the data we're going to store: a list of images with a title, alternative text that stores a description of the image for accessibility, and a URL.

```go
type Image struct {
	Title   string
	AltText string
	URL     string
}
```

Then we can initialise some data:

```go
images := []Image{
    {
      Title: "Sunset",
      AltText: "Clouds at sunset",
      URL: "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
    },
    {
      Title: "Mountain",
      AltText: "A mountain at sunset",
      URL: "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
    },
}
```

Now, set up an HTTP server (we did this in the `http-auth` project), and create an endpoint `/images.json` that returns this data in the response:

```console
> curl 'http://localhost:8080/images.json' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Wed, 03 Aug 2022 18:06:34 GMT
Content-Length: 487

[{"Title":"Sunset","AltText":"Clouds at sunset","URL":"https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"},{"Title":"Mountain","AltText":"A mountain at sunset","URL":"https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"}]
```

> 💡 See the [prep README.md](../prep/README.md#command-line-examples) for an explanation of this command line example.

We're going to need to import `"encoding/json"` and `Marshal` to turn the data into JSON.

```go
b, err := json.Marshal(images)
```

We don't always want to use the exact same field names from our programming language as we do in our public JSON. Instead of `Title`, `AltText`, and `URL`, let's make it so that when we serialise the fields, the keys are in [snake_case](https://en.wikipedia.org/wiki/Snake_case) not [CamelCase](https://en.wikipedia.org/wiki/Camel_case), so that instead our `curl` will behave like this:

```console
> curl 'http://localhost:8080/images.json' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Wed, 03 Aug 2022 18:06:34 GMT
Content-Length: 487

[{"title":"Sunset","alt_text":"Clouds at sunset","url":"https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"},{"title":"Mountain","alt_text":"A mountain at sunset","url":"https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"}]
```

Make sure not to change the field names of the struct when we do so!

Next, add a query parameter `indent` which uses `MarshalIndent` instead (https://pkg.go.dev/encoding/json#MarshalIndent) such that the following snippet works.

```console
> curl 'http://localhost:8080/images.json?indent=2' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Mon, 08 Aug 2022 19:57:51 GMT
Content-Length: 536

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
  }
]
```

The indent value should increase the amount of indentation: `?indent=4` should have 4 spaces, but `?indent=2` should have 2.

To do this, we'll need to investigate the `strconv` and `strings` packages in the Go standard library.

We've now got a working server that responds to requests with data in JSON format, and can format it to be more readable by developers exploring the API.

---

Next we're going to set up a database to store data that our server will use.

We'll use [Postgres](https://www.postgresql.org/), which is an open source relational database. Don't worry if that doesn't mean anything right now. Read the Postgres website to find out the core ideas.

First, install Postgres. We may have been provided with access to Amazon Web Services, which provides Postgres for us via their Relational Database Service. If not, we can also run it on our computer: follow the [instructions on the Postgres website](https://www.postgresql.org/download/).

Our goal is to have a database running that we can connect to using a connection string, which will look something like this: `postgres://user:secret@localhost:5432/mydatabasename`

For easy demoing, we'll assume you have Postgres running locally and connect with `postgresql://localhost`.

Next install [pgAdmin](https://www.pgadmin.org/) (if you haven't already). This is a useful user interface that we'll use to set up the Postgres database.

Open up pgAdmin and add a server that connects to your instance of Postgres.

Then add a database by right-clicking on Databases.

![](./readme-assets/create-database.png)

The SQL generated (see the SQL tab) should read as follows:

```sql
CREATE DATABASE "go-server-database"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;
```

> SQL is Structured Query Language: it's pronounced either "sequel" or Ess-Queue-Ell. Both work. It's used for interacting with various data storage systems. You can find out more about it and learn how to use it [on CodeAcademy](https://www.codecademy.com/learn/learn-sql).

---

Data in Postgres is arranged in tables with columns, like a spreadsheet.

Next create a table that will store our image data. Within the `go-server-database` database, open `schemas`, `public`, and then create a new table.

![](./readme-assets/create-table.png)

Add four columns:

- `id`: type `serial` with "Primary key?" and "Not null?" turned on
- `title`: type `text` with "Not null?" turned on
- `url`: type `text` with "Not null?" turned on
- `alt_text`: type `text`

The SQL tab should look like the following:

```sql
CREATE TABLE public.images
(
    id serial NOT NULL,
    title text NOT NULL,
    url text NOT NULL,
    alt_text text,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.images
    OWNER to postgres;
```

Once the table is created, we can insert some data.

Right click on the table and open `Scripts > INSERT script`.

![](./readme-assets/insert-script.png)

This will open the Query tool, which will allow us to add some data, and give some basic SQL to start with.

```sql
INSERT INTO public.images(
	id, title, url, alt_text)
	VALUES (?, ?, ?, ?);
```

If we run this ("play" button at the top) we will get an error, because we haven't provided any data.

Update the SQL to look like this. We don't need to specify an ID: Postgres will create this.

```sql
INSERT INTO public.images(title, url, alt_text)
	VALUES ('Sunset', 'https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80', 'Clouds at sunset');
```

This should insert some data:

```
INSERT 0 1

Query returned successfully in 132 msec.
```

We can check by deleting this query and running a different one:

```sql
SELECT * from images;
```

We should see the row we added, and a value in the ID column.

Now write a new `INSERT` query with the other image file from our code.

When this is done, the `SELECT` query should return two rows with different IDs and data. Nice.

---

We can now connect our server to the database.

The server and the database are two separate systems that communicate over a network, like our client (`curl`) and server did before. This time, the server will be acting as the "client" by making requests to the database, and the database will responding with data it has stored.

We will be sending requests from our server to the database using SQL: from now on we'll call them queries. At the most basic level, queries can read, write, update and delete data.

Let's get going.

Install pgx: `go get github.com/jackc/pgx/v4`

`pgx` is the library we will use to connect to and query the database.

The first step is to connect. Our server needs to know which database to connect to, and we'll do that by running the server like this:

```
> DATABASE_URL='postgres://localhost:5432/go-server-database' go run .
```

`DATABASE_URL` is an environment variable, and our server can access the value of this variable like this:

```go
os.Getenv("DATABASE_URL")
```

Before going any further, write some code that checks if `DATABASE_URL` is set. If it's not: write a helpful error to `os.Stderr` and exit the process with a non-zero code, to indicate failure (`os.Exit(1)`).

Next, create a connection using `pgx.Connect`. We will need to supply a `context`, which can simply be `context.Background()`. We can [find out about contexts in the Go documentation](https://pkg.go.dev/context), but they are not an important concept for now.

Remember to handle errors in connecting to the database gracefully: handle the error and output a useful error message. We can check this is working by turning off Postgres and starting our server:

```
> DATABASE_URL='postgres://localhost:5432/go-server-database' go run .
unable to connect to database: failed to connect to `host=localhost user=tom database=go-server-database`: dial error (dial tcp 127.0.0.1:5432: connect: connection refused)
exit status 1
```

Once connected to the database, we can look in pgAdmin at the Dashboard for our server: it will show an active session.

It's important to gracefully close our connection to the database, so we need to add `defer conn.Close(context.Background())` immediately after we successfully connect. The `defer` keyword means that this line of code is delayed until the enclosing function returns. [Defer is really useful!](https://go.dev/blog/defer-panic-and-recover)

Now we're connected we can fetch some data.

To get data we'll use `conn.Query`, which is a [method](https://gobyexample.com/methods) defined for the connection (`pgx.Conn`) to the Postgres database.

The SQL will be the same as above, to fetch all the images:

```sql
"SELECT title, url, alt_text FROM public.images"
```

Remember to handle the error from `conn.Query` gracefully.

`conn.Query` returns some `Rows` which we can use to `Scan` for the actual data:

We need a container [slice](https://gobyexample.com/slices) for the image: `var images []Image`

We can iterate over the rows: `for rows.Next() { ... }`

Using `Scan` to extract the data from each row means [passing a pointer to the variables](https://www.digitalocean.com/community/conceptual_articles/understanding-pointers-in-go) you want to fill, so they can be modified:

```go
var title, url, altText string
err = rows.Scan(&title, &url, &altText)
```

After handling the possible error, we can add an image to the list:

```go
images = append(images, Image{Title: title, URL: url, AltText: altText})
```

And that's it! We've not got all images from the database, and our server can now return them as JSON.

The structure of our overall solution now should look something like:

```go
func main() {
    // Check that DATABASE_URL is set

    // Connect to the database

    // Handle a possible connection error

    // Defer closing the connection to when main function exits

    // Fetch images from the database

    // Send a query to the database, returning raw rows

    // Handle query errors

    // Create slice to contain the images

    // Iterate through each row to extract the data
        // Extract the data, passing pointers so the data can be updated in place
        // Append this as a new Image to the images slice

    // Create the handler function that serves the images JSON
        // Handle the indent query parameter
        // Convert images to a byte-array for writing back in a response
        // Write a `Content-Type` header
        // Write the data back to the client

    // Listen & serve on port 8080
}
```

Don't forget error handling!

---

There's now quite a bit of logic in this `main` function and it's getting a bit big.

Extract the image fetching logic (not the database connection) to its own function with this signature:

```go
func fetchImages(conn *pgx.Conn) ([]Image, error) {
    // Fetch images from the database

    // Send a query to the database, returning raw rows

    // Handle query errors

    // Create slice to contain the images

    // Iterate through each row to extract the data
        // Extract the data, passing pointers so the data can be updated in place
        // Append this as a new Image to the images slice

    // Return the images
}
```

---

Now we're going to accept data over an [HTTP POST request](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/POST). This will add new images to the database.

This is going to be an exercise for you. At the end, the following request should work:

```console
> curl 'http://localhost:8080/images.json?indent=2' -i --data '{"title": "Cat", "alt_text": "A cool cat", "url": "https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"}'
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

After that request, a `GET /images.json` request should return all the images:

```console
> curl 'http://localhost:8080/images.json?indent=2' -i
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

## Extensions

1. Don't let the same image URL be uploaded twice
1. On upload, try to detect if completely useless alt-text has been supplied, and reject the request
1. On upload, have our server make sure the URL actually works and returns an image, and reject the request if it didn't
1. On upload, record the resolution of the fetched image in the database too, and return it in the JSON
