.Phony: prerequesits
prerequesits:
## Install migrate
## Install Terraform

.Phony: run-local-sql-migration
run-sql-migration: prerequesits
	migrate -path=./migrations -database=$DB_STRING up