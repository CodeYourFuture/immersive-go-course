# Batch Processing

In this project, you'll build a simple image processing pipeline: reading from a CSV file, downloading an image, processing it, and uploading it to cloud storage.

> ⚠️ This project requires you to have access to an Amazon AWS account, with permissions to configure IAM and S3. Ask on CYF Slack for help with that.

> ⚠️ You need a GitHub repo to complete this because we'll be using [GitHub Actions](https://docs.github.com/en/actions).

Learning objectives:

- What is batch processing and how does it differ from building servers?
- How do we build resilience into batch processing software?
- How do we use Go to use run existing software to complete tasks?
- What is cloud storage technology and how do we read from it & upload data to it?
<!-- * How do we deploy batch processing tasks in the cloud? -->

Timebox: 2 days

## Project

This project starts with some scaffolding. We'll be using a tool called [ImageMagick](https://imagemagick.org/) to process the images. The specific details of it are not important, so this is provided for us.

ImageMagick is software written in C. We're going to call out to this code from Go using a package called [imagick](https://github.com/gographics/imagick), which acts an an interface between our Go code and the ImageMagick code.

As a result of this dependency on ImageMagick, we're going to run this whole project in Docker. That means, rather than use `go run` or `go test`, we'll build and run a Docker container that has ImageMagick correctly installed.

> It's possible to install ImageMagick on your computer to run this, but you don't _need_ to do that.

The Docker commands get a bit complex so we can do things like run tests and _see_ the output of the code, so we're not going to type the Docker commands ourselves.

Instead, to run the code, we use a tool called [Make](https://www.gnu.org/software/make/). Make enables us to build code without knowing the details of how that is done, because these details are recorded in the `Makefile`. In that way it's similar to Docker: the complexity of how to run the code in a container is supplied in the `Dockerfile`.

Read a [quick introduction to Make](https://www.gnu.org/software/make/) and then install it for your computer. On macOS, it's likely that you have it already — try typing `make` into your termincal. If not, you need to [install the command-line tools](https://www.freecodecamp.org/news/install-xcode-command-line-tools/).

> **Important**: you don't need to understand `make` in depth. We're only going to use the very simplest features.

### Run the scaffolding code

Now we can run the code:

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

Here we can see the docker image building. The first time we run this, it will take a while and we won't see lots of `CACHED` messaged. The time is spent installing ImageMagick into the Docker container, and building the app for the first time.

Should now have a new directory — `outputs` — which contains an image `gradient_bw.jpg` which is a grayscale version of the `gradient.jpg` file that's in the `inputs` directory. Something interesting has happened here! Our docker container wrote the file back to _our_ filesystem. It was able to do this because our docker command [mounted](https://docs.docker.com/storage/bind-mounts/) the `outputs` directory from our host (your computer) onto the `/outputs` directory inside the docker container. When the go application wrote the processed file to `/outputs/gradient_bw.jpg`, it was actually writing to your files.

This feature allows us to test the code and see the outputs. Otherwise, they'd be stuck inside the Docker container.

To find out more about this, read the [managing data with Docker](https://docs.docker.com/storage/) guide.

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
