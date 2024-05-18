echo "test\t database tables.."
source ../.env
source ./util.sh # import run_test, summary

case1="have mail table"
have_mail_table=$(docker exec kmail_database_server psql $DATABASE_CONNECTION_STRING -c '\dt' | grep mail)
run_test "$case1" "$have_mail_table"

case2="have used uuid table"
have_used_uuid_table=$(docker exec kmail_database_server psql $DATABASE_CONNECTION_STRING -c '\dt' | grep used_uuid)
run_test "$case2" "$have_used_uuid_table"

summary
