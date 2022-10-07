# buggy-app

Debug and fix a buggy application.

Learning objectives:

- How can we quickly read, understand and fix existing application code?
- How do we QA running code by thinking about security, edge cases and performance?
- How do we ensure safe data access with authentication & authorisation?
- What are some common architectures for services in tech companies?

The following may be covered by previous projects:

- How do we run multiple services and dependencies locally?
- How can services interact beyond HTTP?

## Plan

This will be a simple three-service application: API service, auth service and database.

The API service will pull in a client module for the API service, and communicate with it over gRPC.

### Architecture

```
               ┌───────────────────────────────────────┐      ┌─────────────────┐
               │              API Service              │      │       DB        │
               │                                       │      │                 │
               │ ┌────────────┐           ┌─────────┐  │      │                 │
     ┌────┐    │ │            │           │         │  │      │ ┌─────────────┐ │
     │HTTP│    │ │            │           │         │  │      │ │             │ │
─────┴────┴────┼▶│    Auth    │──────────▶│  Notes  │──┼──────┼▶│    Notes    │ │
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
                │ Auth Service │──────────────────────────────┼▶│    Users    │ │
                │              │                              │ │             │ │
                └──────────────┘                              │ └─────────────┘ │
                                                              └─────────────────┘
```

1. HTTP request hits API service with HTTP simple auth, and auth layer first
1. Auth client calls Auth Service over gRPC, which verifies credentials against Users table
1. Once validated, auth client allows request to continue to the Notes module, which returns the data

Possibly there would be a simple frontend that interacts with the API.

## API

- `GET /1/my/notes` -- Get all notes for the authenticated user
- `GET /1/my/notes/:id` -- Get a specific note for the authenticated user

The Notes model will return a "tags" field with the content of the note, by looking for `#hashtag`.

## DB

The database will be Postgres.

User:

- `id`: primary key
- `status`: `active | inactive`
- `password`: bcrypt
- `created`: timestamp

Note:

- `id`: primary key
- `owner`: fkey into User
- `content`: text
- `created`: timestamp

## File structure

This will follow a mono-repo structure, with `api` and `auth` at the root.

The `auth` package will expose a client that the `api` will depend on.

The services will be coordinated locally with docker-compose.

## Challenge

This application will contain some bugs, as follows:

> ⚠️ Don't read this if you are working through the project!

<details><summary>Bugs</summary>
<p>

1. The Notes model will allow access to any note regardless of the authenticated user
1. The Notes model will query _all_ Notes, and then filter them by ID on the server-side
1. The Notes DB tables will be missing an index on owner
1. The Auth Client will not check if the user is active or inactive
1. The Auth Client will cache authentication results in-memory with no TTL
1. The Notes tags implementation will have a buggy regex that is too eager (`#this will will be a tag up to punctuation`)

</p>
</details>

The instructions will be that there are at least N bugs, and the use the application to find and fix them.
