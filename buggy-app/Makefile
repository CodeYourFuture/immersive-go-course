.PHONY: protoc volumes volumes-reset test run

##
## SETUP
##

volumes/secrets/postgres-passwd:
	mkdir -p volumes/secrets
	# Create a random password for Postgres
	openssl rand -hex 24 | tr -d '\n' > volumes/secrets/postgres-passwd

volumes: volumes/secrets/postgres-passwd
	mkdir -p volumes/data

# Run this to completely reset the database state
# Be careful not to run this while the containers are runing
volumes-reset:
	rm -rf volumes/data volumes/secrets

##
## BUILD
##

# Compile the protobuf files
protoc: auth/service/auth.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		auth/service/auth.proto

##
## TEST
##

test: protoc
	go test -v ./...

##
## RUN
##

run: protoc volumes
	docker-compose up --build
