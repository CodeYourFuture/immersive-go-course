# Batch Processing

In this project, you'll build a simple image processing pipeline: reading a list of image URLs from a CSV file, downloading each image, processing it, and uploading it to cloud storage.

> ⚠️ This project requires you to have access to an Amazon AWS account, with permissions to configure IAM and S3. Ask on CYF Slack for help with that.

> ⚠️ You need a GitHub repo to complete this because we'll be using [GitHub Actions](https://docs.github.com/en/actions).

Learning objectives:

- What is batch processing and how does it differ from building servers?
- How do we use Go to use run existing software to complete tasks?
- What is cloud storage technology and how do we read from it & upload data to it?
- How do we read, modify and extend existing code?
  <!-- - How do we build resilience into batch processing software? -->
  <!-- - How do we deploy batch processing tasks in the cloud? -->

Timebox: 2 days

## Project

This project starts with some scaffolding code already written. We'll be using a tool called [ImageMagick](https://imagemagick.org/) to process the images. The specific details of it are not important, so this is provided for us.

ImageMagick is software written in C. We're going to call out to this code from Go using a package called [imagick](https://github.com/gographics/imagick), which acts an an interface between our Go code and the ImageMagick code.

As a result of this dependency on ImageMagick, we're going to run this whole project in Docker. That means, rather than use `go run` or `go test`, we'll build and run a Docker container that has ImageMagick correctly installed.

> It's possible to install ImageMagick on your computer to run this outside of docker, but you don't _need_ to do that.

The Docker commands get a bit complex so we can do things like run tests and _see_ the output of the code, so we're not going to type the Docker commands ourselves.

Instead, to run the code, we use a tool called [Make](https://www.gnu.org/software/make/). Make enables us to build code without knowing the details of how that is done, because these details are recorded in the `Makefile`. In that way it's similar to Docker - when using Make you store the complexity of what to run in a `Makefile`, and when using Docker you store the complexity of how to run the code in a container in a `Dockerfile`.

Read a [quick introduction to Make](https://www.gnu.org/software/make/) and then install it for your computer. On macOS, it's likely that you have it already — try typing `make` into your terminal. It should output something like `make: no target specified`. If not (e.g. it says something like `command not found: make`), you need to [install the command-line tools](https://www.freecodecamp.org/news/install-xcode-command-line-tools/).

> **Important**: you don't need to understand `make` in depth. We're only going to use the very simplest features. In particular, you shouldn't have to edit the `Makefile` we've supplied other than exactly as we describe (but feel free if you want to!)

### Run the scaffolding code

Now we can run the code, using `make run`. This command looks in the `Makefile` for a "target" called `run`, and then executes the code it finds there. In our case, that looks something like this:

```Makefile
run:
	docker build --target run -t run .
	docker run --rm run
```

We build a docker image and tag it `run`, and then run it.

Let's do that:

```console
> make run
mkdir -p outputs
docker build --target run -t run .
[+] Building 1.1s (21/21) FINISHED
 => [internal] load build definition from Dockerfile                                                 0.0s
 => => transferring dockerfile: 37B                                                                  0.0s
 => [internal] load .dockerignore                                                                    0.0s
 => => transferring context: 2B                                                                      0.0s
 => resolve image config for docker.io/docker/dockerfile:1                                           0.5s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:9ba7531bd80fb0a858632727cf7a112fbfd1  0.0s
 => [internal] load .dockerignore                                                                    0.0s
 => [internal] load build definition from Dockerfile                                                 0.0s
 => [internal] load metadata for docker.io/library/golang:1.19-bullseye                              0.0s
 => [internal] load build context                                                                    0.0s
 => => transferring context: 716B                                                                    0.0s
 => [base  1/12] FROM docker.io/library/golang:1.19-bullseye                                         0.0s
 => CACHED [base  2/12] RUN apt-get update &&     apt-get install -y wget build-essential pkg-confi  0.0s
 => CACHED [base  3/12] RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev     libgif-dev  0.0s
 => CACHED [base  4/12] RUN cd &&     wget https://github.com/ImageMagick/ImageMagick6/archive/6.9.  0.0s
 => CACHED [base  5/12] WORKDIR /app                                                                 0.0s
 => CACHED [base  6/12] COPY go.mod ./                                                               0.0s
 => CACHED [base  7/12] COPY go.sum ./                                                               0.0s
 => CACHED [base  8/12] RUN go mod download                                                          0.0s
 => CACHED [base  9/12] COPY *.go ./                                                                 0.0s
 => CACHED [base 10/12] COPY inputs /inputs                                                          0.0s
 => CACHED [base 11/12] RUN mkdir -p /outputs                                                        0.0s
 => CACHED [base 12/12] RUN go build -o /out                                                         0.0s
 => exporting to image                                                                               0.0s
 => => exporting layers                                                                              0.0s
 => => writing image sha256:155cb66e4ea4ce0d7ec39c451343130a0852ccb7ad312917ac3ec1b0c5b26aa6         0.0s
 => => naming to docker.io/library/run                                                               0.0s

docker run \
		--mount type=bind,source="$(pwd)/outputs",target=/outputs \
		--rm run
2022/10/04 10:45:31 processing: "/inputs/gradient.jpg" to "/outputs/gradient_bw.jpg"
2022/10/04 10:45:32 processed: "/inputs/gradient.jpg" to "/outputs/gradient_bw.jpg"
```

Here we can see the docker image building. The first time we run this, it will take a while and we won't see lots of `CACHED` messages. The time is spent installing ImageMagick into the Docker container and building the app.

One it's run, we should now have a new local directory — `outputs` — which contains an image `gradient_bw.jpg` which is a grayscale version of the `gradient.jpg` file that's in the `inputs` directory.

Something interesting has happened here! Our docker container wrote the file back to _our_ filesystem. It was able to do this because our docker command [mounted](https://docs.docker.com/storage/bind-mounts/) the `outputs` directory from our host (your computer) onto the `/outputs` directory inside the docker container. When the go application wrote the processed file to `/outputs/gradient_bw.jpg`, it was actually writing to your filesystem, outside of the container.

This mount feature allows us to test the code and see the outputs. Otherwise, they'd be stuck inside the Docker container. To find out more about this, read the [managing data with Docker](https://docs.docker.com/storage/) guide.

### Test the scaffolding code

There's also a `make` command for testing, and there are a few tests for the image manipulation code:

```console
> make test
docker build --target test -t test .
[+] Building 1.7s (21/21) FINISHED
 => [internal] load build definition from Dockerfile                                                 0.0s
 => => transferring dockerfile: 37B                                                                  0.0s
 => [internal] load .dockerignore                                                                    0.0s
 => => transferring context: 2B                                                                      0.0s
 => resolve image config for docker.io/docker/dockerfile:1                                           1.2s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:9ba7531bd80fb0a858632727cf7a112fbfd1  0.0s
 => [internal] load .dockerignore                                                                    0.0s
 => [internal] load build definition from Dockerfile                                                 0.0s
 => [internal] load metadata for docker.io/library/golang:1.19-bullseye                              0.0s
 => [base  1/12] FROM docker.io/library/golang:1.19-bullseye                                         0.0s
 => [internal] load build context                                                                    0.0s
 => => transferring context: 716B                                                                    0.0s
 => CACHED [base  2/12] RUN apt-get update &&     apt-get install -y wget build-essential pkg-confi  0.0s
 => CACHED [base  3/12] RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev     libgif-dev  0.0s
 => CACHED [base  4/12] RUN cd &&     wget https://github.com/ImageMagick/ImageMagick6/archive/6.9.  0.0s
 => CACHED [base  5/12] WORKDIR /app                                                                 0.0s
 => CACHED [base  6/12] COPY go.mod ./                                                               0.0s
 => CACHED [base  7/12] COPY go.sum ./                                                               0.0s
 => CACHED [base  8/12] RUN go mod download                                                          0.0s
 => CACHED [base  9/12] COPY *.go ./                                                                 0.0s
 => CACHED [base 10/12] COPY inputs /inputs                                                          0.0s
 => CACHED [base 11/12] RUN mkdir -p /outputs                                                        0.0s
 => CACHED [base 12/12] RUN go build -o /out                                                         0.0s
 => exporting to image                                                                               0.0s
 => => exporting layers                                                                              0.0s
 => => writing image sha256:7f24030e5dc71ed23f98885644321c76b3d6709f7140391940ff06f486a33372         0.0s
 => => naming to docker.io/library/test                                                              0.0s

docker run \
		--rm test
=== RUN   TestGrayscaleMockError
--- PASS: TestGrayscaleMockError (0.00s)
=== RUN   TestGrayscaleMockCall
--- PASS: TestGrayscaleMockCall (0.00s)
PASS
ok  	github.com/CodeYourFuture/immersive-go-course/batch-processing	0.011s
```

This will have run much faster, because all the steps are cached by Docker. This is called the [build cache](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#leverage-build-cache).

Neat! We've now got the application and the tests running in Docker.

## Specificiation

So now we've got some running code, let's get started.

We want a **CLI tool** that:

- reads an input CSV containing URLs
- downloads images by URL
- processes them using ImageMagick to make them monochrome/grayscale
- uploads the results to [Amazon AWS S3 cloud storage](https://aws.amazon.com/s3/)
- writes an output CSV describing what it did

It should run like this: `go run . --input input.csv --output output.csv`

You should modify and extend the `main.go` file we've supplied in this directory, re-using some bits of it.

An example input CSV file would look like this:

```csv
url
https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80
https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80
https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80
```

Or as a table:

| url                                                                                                                                                            |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80 |
| https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80 |
| https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80    |

An example output CSV would look like this:

```csv
url,input,output,s3url
https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80,/tmp/1664885100393-5577006791947779410.jpg,/tmp/1664885100394-8674665223082153551.jpg,https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885100394-8674665223082153551.jpg
https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80,/tmp/1664885101726-6129484611666145821.jpg,/tmp/1664885101726-4037200794235010051.jpg,https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885101726-4037200794235010051.jpg
https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80,/tmp/1664885102225-3916589616287113937.jpg,/tmp/1664885102225-6334824724549167320.jpg,https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885102225-6334824724549167320.jpg
```

Or as a table:

| url                                                                                                                                                            | input                                      | output                                     | s3url                                                                                                         |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------- |
| https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80 | /tmp/1664885100393-5577006791947779410.jpg | /tmp/1664885100394-8674665223082153551.jpg | https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885100394-8674665223082153551.jpg |
| https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80 | /tmp/1664885101726-6129484611666145821.jpg | /tmp/1664885101726-4037200794235010051.jpg | https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885101726-4037200794235010051.jpg |
| https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80    | /tmp/1664885102225-3916589616287113937.jpg | /tmp/1664885102225-6334824724549167320.jpg | https://immersive-go-course-batch-processing.s3.eu-west-1.amazonaws.com/1664885102225-6334824724549167320.jpg |

The tool should (in order of priority):

1. Write thorough logs (`log.Println` and `log.Printf`) to describe what it is doing, including errors
1. Validate the input CSV to ensure it only has **one** column, `url`
1. Gracefully handle failures & continue to process the input CSV even if one row fails
1. Support a configurable AWS region and S3 bucket via environment variables `AWS_REGION` and `S3_BUCKET`

### Project extension

To take this project further, add these requirements:

1. The tool should write an output CSV with the URLs that failed, in a format that could be used as input: `go run . --input input.csv --output output.csv --output-failed failed.csv`
1. Failures in downloading and uploading the images can be temporary: introduce a retry-with-backoff strategy to make the tool more resilient to temporary failures
1. Do not re-upload the same image to S3
   - Can you do this without storing anything in a database?
1. Do not download & process the same image
   - Can you also do this without storing anything in a database?
1. To speed up the tool, process and upload in parallel using [goroutines](https://go.dev/tour/concurrency/1)

## How-to

Most of getting this project build is up to you. However, here are some pointers for things you are going to need.

### Reading a CSV

The built-in [`encoding/csv`](https://pkg.go.dev/encoding/csv) package is the one to use to read and write the CSV files.

### Downloading the file

We can use the the standard `http` package to download the image. Things to watch out for:

- HTTP requests can fail - remember to catch the error!
- HTTP requests can "succeed" but with a non-200 status code. Think about what that could mean!
- How can you make sure the downloaded data is an image, and not some other of file?

### Image processing

Most of the ImageMagick code, which grayscales the image, is written for you. This shouldn't change too much.

### Cloud storage (S3)

The specification asks you to upload the images to [Amazon AWS S3 cloud storage](https://aws.amazon.com/s3/). S3 is a cloud service designed to store huge amounts of data and files in a secure, scalable and cost-effective way.

The basic organisation of S3 is simple. At the top level of S3 is the **bucket**. Data ("objects") are stored within these buckets in folders, like a normal file system. Read the [high-level guide to S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html).

You will need to create a bucket — call it anything you like. Pay attention to which [AWS region](https://cloudacademy.com/blog/aws-regions-and-availability-zones-the-simplest-explanation-you-will-ever-find-around/#:~:text=What%20are%20AWS%20Regions%3F,host%20their%20cloud%20infrastructure%20there.) you create it in – you will need this later.

You can make it publically accessible to the internet, so that you or anyone can view the images that are uploaded, with a policy like this:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicAccess",
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject", "s3:GetObjectVersion"],
      "Resource": "arn:aws:s3:::[NAME OF BUCKET]/*"
    }
  ]
}
```

> ⚠️ Public means public. **DO NOT** upload anything sensitive to the bucket.

Use the [AWS SDK for Go](https://github.com/aws/aws-sdk-go) to interact with Amazon S3.

You will need to set up credentials so that your docker image can write to S3. The best way to do that is to:

- Create a role + policy to allow uploads to the bucket
- Allow your AWS user account to use the role
- Make your AWS credentials available to the Docker container

#### Role + policy

[Create a Policy](https://us-east-1.console.aws.amazon.com/iamv2/home#/policies) `S3ReadWriteGoCourse` with the following configuration, which allows reading and writing to the bucket.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ListObjectsInBucket",
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::[NAME OF BUCKET]"]
    },
    {
      "Sid": "AllObjectActions",
      "Effect": "Allow",
      "Action": "s3:*Object",
      "Resource": ["arn:aws:s3:::[NAME OF BUCKET]/*"]
    }
  ]
}
```

[Create a Role](https://us-east-1.console.aws.amazon.com/iamv2/home#/roles) `GoCourseLambdaUserReadWriteS3` associated with the `S3ReadWriteGoCourse` Policy. Set up the Trust Relationships of this policy as follows:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::[YOUR AWS USER ID]:root"
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

Once that's done, you need local credentials:

1. Install the [AWS CLI](https://aws.amazon.com/cli/)
1. [Follow this guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) to set up your credentials locally

Then change the `Makefile` to [mount](https://docs.docker.com/storage/bind-mounts/) your `.aws` directory into the Docker container, in the `root` user's home directory (`/root/.aws`). When using the Go SDK, it will automatically find these credentials and use them:

```Makefile
run: outputs
	docker ...
        ...
        --mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws \
		...
```

You'll need to say which AWS region the S3 bucket is located in. Do that with the `AWS_REGION` environment variable. You can supply environment variables to the Docker commands directly in the `Makefile`...

```Makefile
run: outputs
	docker ...
        ...
        -e AWS_REGION=eu-west-1 \
		...
```

... or use an [env file](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file):

```Makefile
run: outputs
	docker ...
        ...
		--env-file docker_env \
        ...
```

Here's an example `docker_env` file:

```env
AWS_REGION=eu-west-1
```

### Session

We'll need to set up an AWS session using your Role. To do that, we give it the [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of the Role, which you can from the Role's page in the [AWS IAM dashboard](https://us-east-1.console.aws.amazon.com/iamv2/home?region=us-east-1#/roles). It's best to also pass this to the Docker container via the environment, as we did with the `AWS_REGION`:

```go
// Get the Role ARN
awsRoleArn := ???

// Set up S3 session
sess := ???

// Create the credentials from AssumeRoleProvider to assume the role
// referenced by the ARN.
creds := stscreds.NewCredentials(sess, awsRoleArn)

// Create service client value configured for credentials
// from assumed role.
svc := s3.New(sess, &aws.Config{Credentials: creds})
```

### Uploading

The example code on the [AWS SDK for Go](https://github.com/aws/aws-sdk-go) README will be helpful!
