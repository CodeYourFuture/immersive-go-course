.PHONY: protoc

protoc: auth/service/auth.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		auth/service/auth.proto

test: protoc
	go test -v ./auth