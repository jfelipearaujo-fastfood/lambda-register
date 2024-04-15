#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS customers (id varchar(255), document_id varchar(255), document_type int, is_anonymous boolean, password varchar(255), created_at TIMESTAMP, updated_at TIMESTAMP, PRIMARY KEY (id));
EOSQL
