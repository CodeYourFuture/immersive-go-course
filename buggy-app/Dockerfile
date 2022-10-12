# syntax=docker/dockerfile:1

##
## BASE
##
FROM golang:1.19-bullseye as base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p /out
RUN go build -o /out ./...

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=base /out /out

EXPOSE 8080