# buggy-app

See https://github.com/CodeYourFuture/immersive-go-course/issues/53 for context.

## Implementation

Started with gRPC for auth service, as this is the newest aspect. Very simple interface to start with: `Verify(Input) (Result, error)`

Went down a big rabbit hole with [`Context`](https://go.dev/blog/context) -- hadn't spent much time on this but it turns out to be important. Something we should teach.

### Database

Using the [postgres docker image](https://github.com/docker-library/docs/blob/master/postgres/README.md).

Configuration:

- For running locally, we want to mount a directory on the host machine
- Startup order will be important: https://docs.docker.com/compose/startup-order/
- We need to generate secrets for the postgres user

For migrations, `golang-migrate`:

- `brew install golang-migrate` for the global executable
- We'll run them directly from go: `https://pkg.go.dev/github.com/golang-migrate/migrate`
