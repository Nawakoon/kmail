echo "test\t setup database.."
source ../.env
source ./util.sh # import run_test, summary

database_admin="postgresql://kmail_database_admin:$KMAIL_DB_ADMIN_PASSWORD@localhost:5432/$DATABASE_NAME"

case1="have local database server container"
haveDatabaseContainer=$(docker ps -a | grep kmail_database_server)
run_test "$case1" "$haveDatabaseContainer"

case2="have database for kmail"
haveDatabase=$(docker exec kmail_database_server psql -U postgres -c '\l' | grep $DATABASE_NAME)
run_test "$case2" "$haveDatabase"

case3="database have user for database admin"
databaseHaveAdminUser=$(docker exec kmail_database_server psql -U postgres -c '\du' | grep kmail_database_admin)
run_test "$case3" "$databaseHaveAdminUser"

case4="database admin user can access kmail database"
havePostgresShell=$(which psql)
if [ "$havePostgresShell" ]; then
  databaseAdminCanAccessDatabase=$(psql "$database_admin" -c "\du" | grep kmail_database_admin)
  run_test "$case4" "$databaseAdminCanAccessDatabase"
else
  skip_test "$case4" "psql not found"
fi

summary
