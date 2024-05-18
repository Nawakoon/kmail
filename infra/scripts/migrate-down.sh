#!/bin/bash

echo "$(date): Migration down..."
source ../.env

haveDatabaseContainer=$(docker ps -a | grep kmail_database_server)
if [ ! "$haveDatabaseContainer" ]; then
  echo "kmail database container not found, please setup the database container"
  exit 1
fi

# docker read local migration files to run
docker exec -it kmail_database_server mkdir -p /tmp/migrations

for file in ../migration/down/*.down.sql; do
  docker cp $file kmail_database_server:/tmp/migrations/
  basefile=$(basename $file)
  docker exec -it kmail_database_server psql $DATABASE_CONNECTION_STRING -f /tmp/migrations/$basefile
done

echo "$(date): Migration up complete"