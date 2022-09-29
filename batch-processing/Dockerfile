# syntax=docker/dockerfile:1
## Build
FROM golang:1.19-bullseye as build

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /out

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /out /out

ENTRYPOINT ["/out"]