# buggy-app

See https://github.com/CodeYourFuture/immersive-go-course/issues/53 for context.

## Implementation

Started with gRPC for auth service, as this is the newest aspect.

Very simple interface to start with: `Verify(Input) (Result, error)`

Went down a big rabbit holw with [`Context`](https://go.dev/blog/context) -- hadn't spent much time on this but it turns out to be important. Something we should teach.
