# Docker & Cloud Deployment

In this project, you'll build a simple Go server application and Dockerise it to run within a container. You'll write tests that run against the container, and then build GitHub actions automate continuous testing and release of the application. You'll then run it in the cloud by pushing the container image to Amazon AWS Elastic Container Repository, and then launch it in Amazon AWS Elastic Container Service using the user interface.

> ⚠️ This project requires you to have access to an Amazon AWS account, with permissions to configure ECS, ECR, Fargate, and Elastic Load Balancing. Ask on CYF Slack for help with that.

Learning objectives:

- Det up Docker locally
- `Dockerfile` syntax & how to containerise an application
- Run applications locally using Docker
- Push container images to a repository ([ECR](https://aws.amazon.com/ecr/))
- Deploy images to Elastic Cloud resources using Elastic Container Service
- Destroy cloud resources to mitigate cost

Timebox: 2 days

## Project

This project will require us to pull together information from several different guides to get an application developed locally to run in the cloud. This project will give you experience of building and deploying a real applications. Our focus here is **not** on the Go code, but on the infrastructure around it.

> **Important:** You will need a GitHub repo to complete this becase we'll be using [GitHub Actions](https://docs.github.com/en/actions). If you are not already working in the `immersive-go-course` repository, now would be a good time!

### Motivation

We're going to build a simple Go server, and then _containerise_ it.

> In this project, we'll use the terms containerise and dockerise interchangeably to mean the same thing: making an application run in a container using Docker.

A container is a "sandboxed" process on your computer that is isolated from all other processes, unless specifically allowed.

When running a container, it uses an isolated filesystem. This filesystem is provided by a container image. Since the image contains the container’s filesystem, it must contain everything needed to run an application - all dependencies, configurations, scripts, binaries, etc. The image also contains other configuration for the container, such as environment variables, a default command to run, and other metadata.

When you combine images and containers, we can package whole applications in way that is transferrable (we can create them completely separately from running them) and highly reproducible. Both of these are very important in a production environment.

The isolation and security allows you to run many containers simultaneously on a given host. Containers are lightweight and contain everything needed to run the application, so you do not need to rely on what is currently installed on the host. You can easily share containers while you work, and be sure that everyone you share with gets the same container that works in the same way.

To summarize, a container:

- is a runnable instance of an image. You can create, start, stop, move, or delete containers
- can be run on local machines, virtual machines or deployed to the cloud
- is portable (can be run on any OS)
- is isolated from other containers and runs its own software, binaries, and configurations

## Background

### Docker

Docker is an open platform for developing, shipping, and running applications, based around containers and images. Docker provides tooling and a platform to manage the lifecycle of your containers:

- Develop your application and its supporting components using containers
- The container becomes the unit for distributing and testing your application.
- When you’re ready, deploy your application into your production environment, as a container or an orchestrated service. This works the same whether your production environment is a local data center, a cloud provider, or a hybrid of the two.

Read [this guide to get an overview of Docker](https://docs.docker.com/get-started/overview/).

To build hands-on familiarity with Docker, complete parts 1, 2 and 3 of [this tutorial](https://docs.docker.com/get-started/), after which you should know about:

- Running applications with docker: `docker run -dp 80:80 docker/getting-started`
- Containers and images: process & filesystem isolation
- `Dockerfile`: a text-based script of instructions that is used to create a container image
- Starting, managing processes and images: `docker ps` and `docker rm -f`s

Next work through the [Go Docker tutorial](https://docs.docker.com/language/golang/), after which you should know about:

- Dockerising a go application
- Starting and stopping containers
- Volumes & networking between docker containers
- Basics of docker-compose and CockroachDB
- GitHub actions for pushing the image to Docker Hub

Spend some time on these steps, and feel free to complete other tutorials too. It's very important to grasp the core ideas of containers, images and docker:

- Docker is a set of tools for managing containers and images
- Images are frozen file systems that hold everything a container needs to run
- Containers are your running application, based on an image

### Cloud hosting in AWS

We're going to host our application in the cloud, specifically in Amazon Web Services (AWS). AWS is a large suite of products for running technology in the cloud.

The set of AWS products we're going to interact with directly are:

- Elastic Container Repository (ECS): store images that can later be run as containers
- Elastic Container Service (ECR): run containers, including all the infrastructure needed to make them accessible to the internet
- Identity & Access Management (IAM): manage security, identity and access within AWS

To get familiar with ECS, run through the [AWS tutorial](https://aws.amazon.com/getting-started/hands-on/deploy-docker-containers/), after which you should know about:

- Container & task: like a blueprint for your application
- Service & load balancing: launches and maintains copies of the task definition in your cluster
- Cluster: compute resources used to run the service & load balancing

## Building & Dockerising server

The rest of this project will cover putting this all together to run an application that we've written on Elastic Container Service. The steps will be:

- Build a simple Go server
- Dockerise it to run locally within a container
- Write tests that run against the docker container
- Build GitHub actions automate CI/CD
- Push the image to ECR
- Launch it in ECS using the UI

Make sure to commit code as you work with clear messages. We are going to work in a tight loop with GitHub — pushing code and testing it — so [good Git hygiene](https://betterprogramming.pub/six-rules-for-good-git-hygiene-5006cf9e9e2) is important.

## Server

Write a simple server in Go with this behaviour when you `go run` it:

```console
> curl localhost:8090
Hello, world.

> curl localhost:8090/ping
pong
```

Make sure that the port is configurable with an environment variable `HTTP_PORT`, and defaults to port `80`.

### Dockerise

Write a `Dockerfile` for the application. Optionally, make this a bit harder by including a multi-stage build.

It should build & run like this:

```console
> docker build . -t docker-cloud
[+] Building 22.6s (15/15) FINISHED
...

> docker run -dp 8090:80 docker-cloud
306cf309f3970d5380cd07c3a54aead7ee8cf4f6726b752fecaec39e40da69f5

> curl localhost:8090
Hello, world.

> curl localhost:8090/ping
pong
```

### Tests

Write some simple tests for your server. For writing tests, we'll use [dockertest](https://github.com/ory/dockertest).

The principle of `dockertest` is to test against real running services, end to end, in containers. The advantage is that these services can be destroyed after testing and recreated from scratch, so that the tests are highly reproducible.

Following [docs here](https://github.com/ory/dockertest) and [example here](https://github.com/olliefr/docker-gs-ping), write some tests.

Go dependencies are installed by `go get`.

Make sure to `COPY go.sum ./` in your `Dockerfile`. [Read this guide to understand why](https://golangbyexample.com/go-mod-sum-module/).

## CI/CD

A set of good practices has exists for developing software for production use by thousands or millions of people. One example of this is Continuous Integration & Continuous Deployment, referred to as CI/CD. We're going to focus on the "CI" component of this, which means:

- working with version control (Git)
- running automated checks on the code added to Git (specifically, test)
- automating the steps required to "build" a version of the application that could be deployed to production (specifically, building a Docker image and pushing it to a shared repository)

This automation runs in the cloud — _not_ on developer laptops — so that it is highly reproducible and the same for whoever writes the the code. The idea is that, should the tests fail, the code will not be built or pushed to the repository, so it's harder to push buggy code to production.

### Github Action: running tests

The system we'll use for CI testing and image creation is [GitHub Actions](https://docs.github.com/en/actions).

Follow [this guide on GitHub](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go) to get a GitHub action testing the code.

You'll be adding a single file to the `.github/workflows` directory of the repo.

Pay attention to the following:

```yml
# Make sure that the working directory for the tests is correct
defaults:
  run:
    working-directory: docker-cloud
```

And:

```yml
# Use the right go version file and dependency path
# Note that these are *not* subject to the working-directory defaults
- name: Set up Go
  uses: actions/setup-go@v3
  with:
    go-version-file: "docker-cloud/go.mod"
    cache-dependency-path: "docker-cloud/go.sum"
    cache: true
```

### Github Action: publish image

TODO: @tgvashworth finish this section

Next, we're going build an action that creates and pushes the image to AWS Elastic Container Registry.

This is complex.

The [guide here](https://benoitboure.com/securely-access-your-aws-resources-from-github-actions) was very helpful.

Actions used:

- https://github.com/aws-actions/configure-aws-credentials
- https://github.com/aws-actions/amazon-ecr-login

Consider: setting up a CYF public ECR repository with the right permissions.

## Running on ECS

TODO: @tgvashworth write this section
