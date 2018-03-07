#!/bin/sh


echo "L'host:"
echo "$MYSQL_HOST"
echo "Il db:"
echo "$MYSQL_DATABASE"
echo "Lo user:"
echo "$MYSQL_USER"
echo "La password:"
echo "$MYSQL_PASSWORD"
echo "Il file:"
echo "$MYSQL_SQL_FILENAME"

_now=$(date +"%s_%m_%d_%Y")
_file="dumpdb/$MYSQL_SQL_FILENAME_$_now.sql"

echo $_file

mysqldump  "$MYSQL_DATABASE" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -h "$MYSQL_HOST" > $_file