#!/bin/bash

set -e
set -u
set -x
set -o pipefail

export PGPASSWORD="$POSTGRES_PASSWORD"
export PGUSER="$POSTGRES_USER"

GF_DATABASE_NAME="grafana"
for sql in \
        "CREATE USER \"${GF_DATABASE_USER}\" WITH LOGIN PASSWORD '${GF_DATABASE_PASSWORD}' NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT;" \
        "ALTER DATABASE \"${GF_DATABASE_NAME}\" OWNER TO \"${GF_DATABASE_USER}\";" \
        "ALTER ROLE \"${GF_DATABASE_USER}\" SET client_encoding TO \"utf8\";" \
        "ALTER ROLE \"${GF_DATABASE_USER}\" SET default_transaction_isolation TO \"read committed\";" \
        "ALTER ROLE \"${GF_DATABASE_USER}\" SET timezone TO \"UTC\";" \
        "GRANT ALL PRIVILEGES ON DATABASE \"${GF_DATABASE_NAME}\" TO \"${GF_DATABASE_USER}\";" \
        "GRANT CREATE ON DATABASE \"${GF_DATABASE_NAME}\" TO \"${GF_DATABASE_USER}\";" \
        "GRANT ALL ON ALL TABLES IN SCHEMA public TO \"${GF_DATABASE_USER}\";" \
        "GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO \"${GF_DATABASE_USER}\";" \
        "GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO \"${GF_DATABASE_USER}\";" \
; do
  psql "$GF_DATABASE_NAME" -U "$PGUSER" -c "$sql";
done;


for sql in \
        "GRANT CONNECT ON DATABASE \"${APP_DATABASE_NAME}\" TO \"${GF_DATABASE_USER}\";" \
        "GRANT USAGE ON SCHEMA public TO \"${GF_DATABASE_USER}\";" \
        "GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"${GF_DATABASE_USER}\";" \
        "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO \"${GF_DATABASE_USER}\";" \
; do
  psql "$APP_DATABASE_NAME" -U "$PGUSER" -c "$sql";
done;
