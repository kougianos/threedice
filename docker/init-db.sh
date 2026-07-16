#!/bin/bash
# Creates the Go service's database alongside the Java one. POSTGRES_DB already
# creates "threedice"; the two services own separate schemas so each can run its
# own migration tool without fighting over a shared history table.
#
# Postgres only runs this on first initialisation of the data volume. On an
# existing volume, create it by hand instead:
#   docker compose exec postgres psql -U threedice -d threedice \
#     -c 'CREATE DATABASE threedice_go OWNER threedice;'
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE threedice_go OWNER $POSTGRES_USER;
EOSQL

echo "init-db: created database threedice_go"
