# batch-processing

See https://github.com/CodeYourFuture/immersive-go-course/issues/26 for context.

## Plan

The planned architecture of this:

1. Read the CSV
2. Download the images to a location (`/tmp`)
3. Use imagemagick to monochrome them
4. Upload them to S3
5. Return the URL

I tried [an initial implementation](https://github.com/CodeYourFuture/immersive-go-course/pull/46) of this that went a long way, but that I didn't like in the end.

The first step will be to build this linearly, and to write tests as we go. Because there is real file getting and writing, we will run integration tests in Docker:

1. Mock the `jpg` get
2. Write a real file
3. Mock S3 methods using [s3iface](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3iface/)

Then use goroutines to run it in parallel, likely by wrapping the output in a mutex and locking/unlocking as the goroutine completes: https://pkg.go.dev/sync#Mutex

A possible last extension would be to use channels: https://go.dev/blog/pipelines

### Downloads

The download is simple — create a file in a temporary location, and `http.Get` into it with `io.Copy`.

### `imagemagick`

To run ImageMagick (and this whole thing) in a repeatable way, we will do it all in a Docker container based on `dpokidov/imagemagick:latest-bullseye` using multi-stage build. This will give us the `magick` command.
