#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER myuser WITH ENCRYPTED PASSWORD 'mypassword';
    CREATE DATABASE mydatabase;
    GRANT ALL PRIVILEGES ON DATABASE mydatabase TO myuser;
EOSQL