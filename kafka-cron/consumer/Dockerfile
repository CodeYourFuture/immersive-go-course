FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./consumer ./
COPY ./types ./types

RUN go build -o /build