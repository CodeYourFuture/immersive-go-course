# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
RUN go build -o /docker-cloud

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /docker-cloud /docker-cloud
EXPOSE 80

ENTRYPOINT [ "/docker-cloud" ]