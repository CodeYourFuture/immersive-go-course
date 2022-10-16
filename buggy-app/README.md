# Buggy App

In this project, you're going to run, debug, and fix a buggy application. This code here is a "notes" application with users and notes. It simulates a real-world application that has grown and developed over time: the code isn't perfect, there are tests missing and it contains bugs. The task is to find and fix these bugs. There are _at least_ 5 distinct and important bugs for you to find a fix.

> **Note:** This project requires you to have Docker with [Compose](https://docs.docker.com/compose/) installed. Compose is a tool for defining and running multi-container Docker applications. With Compose, we use a YAML file to configure the application’s services. Then, with a single command, we create and start all the services from the configuration.

Learning objectives:

- How can we quickly read, understand and fix existing application code?
- How do we QA running code by thinking about security, edge cases and performance?
- How do we ensure safe data access with authentication & authorisation?
- What are some common architectures for services in tech companies?

## Notes App

There are two services and a database in this project:

1. API service: authenticate and get notes over HTTP
1. Auth service: gRPC service for verifying authentication information
1. Database: Postgres storing users and notes

Here's how requests flow through the architecture diagram below:

1. An HTTP request hits API service with HTTP simple auth
1. The API uses the Auth Client to check the authentication information
1. Auth Client calls Auth Service over gRPC, which verifies credentials against User table in the database
1. Once validated, the API allows request to continue to retrieve Note data from the database

```
               ┌───────────────────────────────────────┐      ┌─────────────────┐
               │              API Service              │      │       DB        │
               │                                       │      │                 │
               │ ┌────────────┐           ┌─────────┐  │      │                 │
     ┌────┐    │ │            │           │         │  │      │ ┌─────────────┐ │
     │HTTP│    │ │            │           │         │  │      │ │             │ │
─────┴────┴────┼▶│    Auth    │──────────▶│  Notes  │──┼──────┼▶│    Note     │ │
               │ │            │           │         │  │      │ │             │ │
               │ │            │           │         │  │      │ └─────────────┘ │
               │ └────────────┘           └─────────┘  │      │                 │
               │        ▲                              │      │                 │
               └────────┼──────────────────────────────┘      │                 │
                        ├────┐                                │                 │
                        │gRPC│                                │                 │
                        ├────┘                                │                 │
                        ▼                                     │                 │
                ┌──────────────┐                              │ ┌─────────────┐ │
                │              │                              │ │             │ │
                │ Auth Service │──────────────────────────────┼▶│    User     │ │
                │              │                              │ │             │ │
                └──────────────┘                              │ └─────────────┘ │
                                                              └─────────────────┘
```

## API

- `GET /1/my/notes.json` -- Get all notes owned by the authenticated user
- `GET /1/my/notes/:id.json` -- Get a specific note owned by the authenticated user

Authentication is by [basic auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication):

```console
> curl 127.0.0.1:8090/1/my/notes.json \
    -H 'Authorization: Basic QTJSUHE2VG86YmFuYW5h' -i
HTTP/1.1 200 OK
Content-Type: text/json
Date: Sun, 16 Oct 2022 09:45:03 GMT
Content-Length: 162

{"notes":[{"id":"JBmytGF3","owner":"A2RPq6To","content":"Example note content #exampletag","created":"2022-10-15T19:48:19.597524Z","modified":"2022-10-15T19:48:19.597524Z", "tags": ["exampletag"]}]}
```

The API exposes the "tags" associated with a Note. These are not stored, but are extracted as notes are read from the database.

## Database

The database is Postgres. This is the table structure:

### `user`

- `id`: primary key: randomly generated string, like `A2RPq6To`
- `status`: 0 or 1 (inactive, active)
- `password`: bcrypt
- `created`: timestamp
- `modified`: timestamp

Users with status 0 should not be able to authenticate or access their notes.

### `note`

- `id`: primary key: randomly generated string, like `JBmytGF3`
- `owner`: foreign key for a user
- `content`: text, contents of the Note
- `created`: timestamp
- `modified`: timestamp

Users should not be able to access notes that they do not own.

## Structure

Here's what each directory contains:

- `api`: The HTTP API service
  - `model`: Code for interacting with notes in the database
- `assets`: Static files relating to the application (e.g. `.monopic` architecture file)
- `auth`: The Auth service that verifies authentication information supplied to the API service, and an Client that the API service uses to talk to the Auth service
  - `cache`: A caching package that stores previously verified authentication information
  - `service`: Protocol Buffer code (`.proto` and generated `.go`) for the gRPC service
- `bin`: Executable scripts that are used within the Dockerfile
- `cmd`: Command line tools for running the application, setting up the database and generating data for testing
  - `api`: Run the API service
  - `auth`: Run the Auth service
  - `migrate`: Set up the database. See [Migrations](#migrations) below.
- `migrations`: `sql` files for the migrations, setting up `user` and `note` tables
- `util`: Shared code across the other directories
- `volumes`: Directories that will be mounted into the containers
  - `init`: [Scripts for initialising the Postgres database](https://github.com/docker-library/docs/blob/master/postgres/README.md#initialization-scripts)
  - `secrets`: Created when the app is run. Contains secrets such as the `postgres` user password.

In addition there are some important files:

- `docker-compose.yml`: Configuration for `docker compose`, specifying how the database, services, migrations and tests run within docker
- `Dockerfile`: Build configuration for the application. The whole repository is built with the same `Dockerfile` configuration, with different commands available (those in the `cmd` directory). See the `docker-compose.yml` file for details on how each service actually runs.
- `Makefile`: Commands for setting-up, testing, migrating, and running the application. The important commands are:
  - `migrate`: Runs migrations against Postgres. May also do initialisation if required.
  - `test`: Runs tests. May also do initialisation if required.
  - `build-run`: Builds and runs the whole application.

## Running

Use `make build run`:

```console
> make build run
...
buggy-app-postgres-1  | 2022-10-16 09:41:48.815 UTC [1] LOG:  database system is ready to accept connections
buggy-app-auth-1      | wait-for-it.sh: postgres:5432 is available after 1 seconds
buggy-app-auth-1      | 2022/10/16 09:41:48 auth service: listening: :80
buggy-app-api-1       | wait-for-it.sh: postgres:5432 is available after 1 seconds
buggy-app-api-1       | 2022/10/16 09:41:49 api service: listening: :80
```

Once it's running, the port of the API will be available (`8090`) which we can see via `docker compose ps`:

```console
> docker compose ps
NAME                   COMMAND                  SERVICE             STATUS              PORTS
buggy-app-api-1        "/bin/docker-entrypo…"   api                 running             127.0.0.1:8090->80/tcp
buggy-app-auth-1       "/bin/docker-entrypo…"   auth                running             127.0.0.1:8080->80/tcp
buggy-app-postgres-1   "docker-entrypoint.s…"   postgres            running             0.0.0.0:5432->5432/tcp
```

Under the hood, we're using `docker compose` to coordinate startup.

`bin/wait-for-it.sh` is used extensively to make sure that Postgres is available before the other services are started.

We can also re-run everything without rebuilding: `make run`

## Tests

To run the tests of this project, run:

```console
> make test
...
ok  	github.com/CodeYourFuture/immersive-go-course/buggy-app/api	2.051s
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/api/model	[no test files]
ok  	github.com/CodeYourFuture/immersive-go-course/buggy-app/auth	5.974s
ok  	github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/cache	0.002s
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/cmd/api	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/cmd/auth	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/cmd/migrate	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/cmd/test	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/util	[no test files]
?   	github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authuserctx	[no test files]
```

**Important:** the tests run **inside Docker** and rely on a fully [migrated](#Migrations) Postgres. Always running via `make` should ensure this is the case.

## Migrations

In this context, database migrations are SQL files (`.sql`) that specify how the data should be setup. They are ordered, and the order is very important: each migration file builds on the previous migrations so that we get a fully working database state at the end. The migration files are used to set up database tables, plus functions and triggers for generating random IDs.

Each migration has an "up" and "down" script which performs the migration and undoes it, respectively.

The migration process is performed by the code in `cmd/migrate` using the [migrate package](https://github.com/golang-migrate/migrate).

To run migrations:

```console
> make migrate
...
2022/10/16 10:17:19 migrate: "file:///migrations/app" into "app" database
...
2022/10/16 10:17:19 migrate: complete
```

**Important:** the migrations run **inside Docker**. Always run them via `make`.
