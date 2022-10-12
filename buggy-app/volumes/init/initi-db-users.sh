#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER buggy-app;
	CREATE DATABASE buggy-app;
	GRANT ALL PRIVILEGES ON DATABASE buggy-app TO buggy-app;
EOSQL