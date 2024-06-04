FROM golang:1.22.3

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY main/*.go ./main/
COPY raft_proto/*.go ./raft_proto/


EXPOSE 7600

# Bit lazy not to build properly but that's not the main point of this exercise
CMD ["go", "run", "main/main.go"]