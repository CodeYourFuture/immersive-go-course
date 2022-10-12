# syntax=docker/dockerfile:1

# This Dockerfile contains all code for the entire repository.
#
# To run a different executable, supply a different command.
# To avoid the "wait for Postgres" feature, supply a different entrypoint.
FROM golang:1.19-bullseye as base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p /out
RUN go build -o /out ./...

COPY bin /bin
COPY migrations /migrations

# The entrypoint will, by default, wait for postgres to become available at `postgres://postgres:5432` before
# running the command that follows
ENTRYPOINT [ "/bin/wait-for-it.sh", "postgres:5432", "-t", "10", "--" ]