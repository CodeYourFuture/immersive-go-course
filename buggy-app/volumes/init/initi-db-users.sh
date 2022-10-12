#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER auth;
	CREATE DATABASE auth;
	GRANT ALL PRIVILEGES ON DATABASE auth TO auth;
EOSQL