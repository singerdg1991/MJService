#!/bin/bash
echo "Creating database: $DATABASE_NAME"
set -e
export PGPASSWORD=$DATABASE_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$DATABASE_USER" --dbname "$DATABASE_DEFAULT_NAME" <<-EOSQL
  CREATE USER $DATABASE_DEFAULT_NAME WITH PASSWORD '$DATABASE_PASSWORD';
  CREATE DATABASE $DATABASE_NAME;
  GRANT ALL PRIVILEGES ON DATABASE $DATABASE_NAME TO $DATABASE_USER;
EOSQL