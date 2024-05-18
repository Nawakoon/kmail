#!/bin/bash

echo "$(date): Setup Postgres database..."
source ../.env

CONTAINER_NAME="kmail_database_server"
START_CONTAINER_TIMEOUT=2

postgresImageExists=$(docker images | grep postgres)
if [ ! "$postgresImageExists" ]; then
  docker pull postgres
else
  echo "Postgres image already setup"
fi

postgresContainerExists=$(docker ps -a | grep $CONTAINER_NAME)
if [ ! "$postgresContainerExists" ]; then
  docker run --name $CONTAINER_NAME -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_DB=kmail_test_database -d -p $POSTGRES_PORT:$POSTGRES_PORT postgres && sleep $START_CONTAINER_TIMEOUT
else
  echo "Postgres container already setup"
fi

postgresContainerRunning=$(docker ps | grep $CONTAINER_NAME)
if [ ! "$postgresContainerRunning" ]; then
  docker start $CONTAINER_NAME && sleep $START_CONTAINER_TIMEOUT
else
  echo "Postgres container already running"
fi

databaseExists=$(docker exec $CONTAINER_NAME psql -U postgres -c '\l' | grep $DATABASE_NAME)
if [ ! "$databaseExists" ]; then
  docker exec -it $CONTAINER_NAME psql -U postgres -c "CREATE DATABASE $DATABASE_NAME;"
else
  echo "Database $DATABASE_NAME already setup"
fi

adminUserExists=$(docker exec $CONTAINER_NAME psql -U postgres -c '\du' | grep kmail_database_admin)
if [ ! "$adminUserExists" ]; then
  docker exec -it $CONTAINER_NAME psql -U postgres -c "CREATE USER kmail_database_admin WITH PASSWORD '$KMAIL_DB_ADMIN_PASSWORD';"
  docker exec -it $CONTAINER_NAME psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DATABASE_NAME TO kmail_database_admin;"
  docker exec -it $CONTAINER_NAME psql -U postgres -c "ALTER USER kmail_database_admin CREATEDB;"
  docker exec -it $CONTAINER_NAME psql -U postgres -c "ALTER USER kmail_database_admin WITH SUPERUSER;"
else
  echo "User kmail_database_admin already setup"
fi

echo "$(date): Setup Postgres database complete"