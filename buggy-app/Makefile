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
	rm -rf volumes/data

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

test: volumes
	docker compose --profile test build
	docker compose --profile test run test

##
## RUN
##

build: volumes
	docker compose --profile run build

migrate: volumes
	docker compose --profile migrate build
	docker compose run migrate

run:
	docker compose --profile run up

build-run: | build migrate run

migrate-local:
	POSTGRES_PASSWORD_FILE=volumes/secrets/postgres-passwd \
		go run ./cmd/migrate --hostport localhost:5432 --path migrations up

migrate-local-down:
	POSTGRES_PASSWORD_FILE=volumes/secrets/postgres-passwd \
		go run ./cmd/migrate --hostport localhost:5432 --path migrations down
