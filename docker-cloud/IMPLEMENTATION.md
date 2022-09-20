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
- Push the image to ECR (not Docker Hub)
- Launch it in ECS using the UI
- Build GitHub actions automate CI/CD
