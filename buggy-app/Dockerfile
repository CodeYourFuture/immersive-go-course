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

COPY bin /bin
COPY migrations /migrations

EXPOSE 8080

ENTRYPOINT [ "/bin/wait-for-it.sh", "postgres:5432", "-t", "10", "--" ]