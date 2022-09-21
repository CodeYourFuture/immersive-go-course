# docker-cloud implementation

To get familiar with Docker, complete parts 1, 2 and 3 of [this tutorial](https://docs.docker.com/get-started/):

- Running applications with docker: `docker run -dp 80:80 docker/getting-started`
- Containers and images: process & filesystem isolation
- `Dockerfile`: a text-based script of instructions that is used to create a container image
- Starting, managing processes and images: `docker ps` and `docker rm -f`s

Next work through the [Go Docker tutorial](https://docs.docker.com/language/golang/):

- Dockerising a go application
- Starting and stopping containers
- Volumes & networking between docker containers
- Basics of docker-compose and CockroachDB
- GitHub actions for pushing the image to Docker Hub

To get familiar with ECS, run through the [AWS tutorial](https://aws.amazon.com/getting-started/hands-on/deploy-docker-containers/):

- Container & task: like a blueprint for your application
- Service & load balancing: launches and maintains copies of the task definition in your cluster
- Cluster: compute resources used to run the service & load balancing

---

The task will be to bring this all together to run a local application on Elastic Container Service:

- Build a simple Go server
- Dockerise it to run locally within a container
- Write tests that run against the docker container
- Build GitHub actions automate CI/CD
- Push the image to ECR (not Docker Hub)
- Launch it in ECS using the UI

## Server

```console
> curl localhost:8090/ping
pong
```

`Dockerfile` inc multi-stage build:

```Dockerfile
# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-bullseye as build

WORKDIR /app

COPY go.mod .
# COPY go.sum .

RUN go mod download

COPY *.go ./

RUN go build -o /out

CMD [ "/out" ]

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /out /out

EXPOSE 80

ENTRYPOINT ["/out"]
```

Build & run:

```console
> docker build . -t docker-cloud
[+] Building 22.6s (15/15) FINISHED
...
> docker run -dp 8090:8090 docker-cloud
306cf309f3970d5380cd07c3a54aead7ee8cf4f6726b752fecaec39e40da69f5
> curl localhost:8090/ping
pong
```

## Tests

[dockertest](https://github.com/ory/dockertest) — principle is to test against real running services, end to end

```console
go get -u github.com/ory/dockertest/v3
```

Following [docs here](https://github.com/ory/dockertest) and [example here](https://github.com/olliefr/docker-gs-ping), write some tests.

## GitHub action to run tests

Follow [guide on GitHub](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go) to get GitHub action testing the code.

Pay attention to:

```yml
defaults:
  run:
    working-directory: docker-cloud
```

And:

```yml
- name: Set up Go
  uses: actions/setup-go@v3
  with:
    go-version-file: "docker-cloud/go.mod"
    cache-dependency-path: "docker-cloud/go.sum"
    cache: true
```

### Action to publish image

This is complex.

[Guide here](https://benoitboure.com/securely-access-your-aws-resources-from-github-actions) was very helpful.

Actions used:

- https://github.com/aws-actions/configure-aws-credentials
- https://github.com/aws-actions/amazon-ecr-login

Consider: setting up a CYF public ECR repository with the right permissions.
