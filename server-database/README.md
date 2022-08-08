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

Write this back as the response:

```
> curl 'http://localhost:8080/images.json' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Wed, 03 Aug 2022 18:06:34 GMT
Content-Length: 487

[{"Title":"Sunset","AltText":"Clouds at sunset","Url":"https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"},{"Title":"Mountain","AltText":"A mountain at sunset","Url":"https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"}]
```

Add a query param `indent` which uses `MarshalIndent` instead (https://pkg.go.dev/encoding/json#MarshalIndent) such that the following snippet works. The indent value should increase the amount of indentation: `?indent=4` should have 4 spaces, but `?indent=2` should have 2.

To do this, you'll need to investigate the `strconv` and `strings` packages in the Go standard library.

```
> curl 'http://localhost:8080/images.json?indent=2' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Mon, 08 Aug 2022 19:57:51 GMT
Content-Length: 536

[
  {
    "Title": "Sunset",
    "AltText": "Clouds at sunset",
    "Url": "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"
  },
  {
    "Title": "Mountain",
    "AltText": "A mountain at sunset",
    "Url": "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\u0026auto=format\u0026fit=crop\u0026w=1000\u0026q=80"
  }
]
```

We've now got a working server that responds to requests with data in JSON format, and can format it.

---

Next we're going to set up a database to store data that our server will use.

We'll use [Postgres](https://www.postgresql.org/), which is an open source relational database. Don't worry if that doesn't mean anything right now. Read the Postgres website to find out the core ideas.

First, install Postgres. You may have been provided with access to Amazon Web Services, which provides Postgres for you via their Relational Database Service. If not, you can also run it on your computer: follow the [instructions on the Postgres website](https://www.postgresql.org/download/).

Your goal is to have a database running that you can connect to using a connection string, which will look something like this: `postgres://user:secret@localhost:5432/mydatabasename`

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

If you run this ("play" button at the top) you will get an error, because we haven't provided any data.

Update the SQL to look like this. We don't need to specify an ID: Postgres will do this.

```sql
INSERT INTO public.images(title, url, alt_text)
	VALUES ('Sunset', 'https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80', 'Clouds at sunset');
```

This should insert some data. We can check by deleting this query and running a different one:

```sql
SELECT * from images;
```

You should see the row you added, and a value in the ID column.

Now write a new `INSERT` query with the other image file from our code.

When this is done, the `SELECT` query should return two rows with different IDs and data. Nice.
