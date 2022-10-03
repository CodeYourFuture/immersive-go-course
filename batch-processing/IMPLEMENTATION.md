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

To be able to run the tests and the app, we end up with multiple targets:

```Dockerfile
FROM golang:1.19-bullseye as base

# ... install dependencies & build ...

FROM base as test

# ... run tests ...

FROM base as run

# ... run app ...
```

Which can then be built by specifying the `--target`:

```console
> docker build --target test -t test .
```

## Developing locally

We can run locally. A few things are needed.

In VSCode settings, if using the go extension:

```json
"gopls": {
    "build.env": {
        "CGO_CFLAGS_ALLOW": "-Xpreprocessor"
    }
}
```

On the CLI:

```console
export PKG_CONFIG_PATH="/usr/local/opt/imagemagick@6/lib/pkgconfig"
```

### Developing in Docker

To develop the app with Docker, we need a slightly fancier command:

```Makefile
develop:
    mkdir -p mount
    docker build --target develop -t develop .
    docker run -it --mount type=bind,source="$$(pwd)",target=/app --mount type=bind,source="/tmp",target=/tmp --rm develop
    rm -rf ./mount
```

## Grayscale

`convert`, accessed via `ConvertImageCommand`, with `-set colorspace Gray -separate -average` seems to work well.

## Upload to S3

- Get credentials set up — https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html
- `brew install awscli`
- `aws configure`

Follow upload example here: `https://github.com/aws/aws-sdk-go`

We need to mount creds from host: `--mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws`

Create `S3ReadWriteGoCourse` policy for IAM role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ListObjectsInBucket",
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::[ID]"]
    },
    {
      "Sid": "AllObjectActions",
      "Effect": "Allow",
      "Action": "s3:*Object",
      "Resource": ["arn:aws:s3:::[ID]/*"]
    }
  ]
}
```

Create `GoCourseLambdaUserReadWriteS3` Role allowing accounts + Lambda to read/write, trust policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::[ID]:root"
      },
      "Action": "sts:AssumeRole"
    },
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

We can then load using URN passed via env:

```go
// Set up S3 session
// All clients require a Session. The Session provides the client with
// shared configuration such as region, endpoint, and credentials.
sess := session.Must(session.NewSession())

// Create the credentials from AssumeRoleProvider to assume the role
// referenced by the ARN.
creds := stscreds.NewCredentials(sess, awsRoleUrn)

// Create service client value configured for credentials
// from assumed role.
svc := s3.New(sess, &aws.Config{Credentials: creds})
```

Need to create a `docker_env` file with config:

```env
AWS_REGION=eu-west-1
AWS_ROLE_URN=arn:aws:iam::[ID]:role/GoCourseLambdaUserReadWriteS3
S3_BUCKET=[ID]
```
