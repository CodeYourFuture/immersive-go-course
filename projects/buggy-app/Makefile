.PHONY: protoc volumes volumes-reset test run

##
## SETUP
##

volumes/secrets/postgres-passwd:
	mkdir -p volumes/secrets
	# Create a random password for Postgres
	openssl rand -hex 24 | tr -d '\n' > volumes/secrets/postgres-passwd

volumes: volumes/secrets/postgres-passwd
	mkdir -p /tmp/buggy-app-data

# Run this to completely reset the database state
# Be careful not to run this while the containers are runing
volumes-reset:
	rm -rf /tmp/buggy-app-data

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
	@# If it doesn't build, we want to know ASAP
	go build ./...
	docker compose --profile test build
	@# -T forces docker not to allocate tty, which is important for pre-commit hooks
	docker compose run -T test

##
## RUN
##

build: volumes
	@# If it doesn't build, we want to know ASAP
	go build ./...
	docker compose --profile run build

migrate: volumes
	docker compose --profile migrate build
	docker compose run migrate

run:
	docker compose --profile run up

run-database:
	docker compose up

build-run: | build run

migrate-local:
	POSTGRES_PASSWORD_FILE=volumes/secrets/postgres-passwd \
		go run ./cmd/migrate --hostport localhost:5432 --path migrations up

migrate-local-down:
	POSTGRES_PASSWORD_FILE=volumes/secrets/postgres-passwd \
		go run ./cmd/migrate --hostport localhost:5432 --path migrations down
