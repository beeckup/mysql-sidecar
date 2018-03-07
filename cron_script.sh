#!/usr/bin/env bash

echo "Il db:"
echo "$MYSQL_DATABASE"
echo "Lo user:"
echo "$MYSQL_USER"
echo "La password:"
echo "$MYSQL_PASSWORD"
echo "Il file:"
echo "$MYSQL_SQL_FILENAME"

_now=$(date +"%m_%d_%Y")
_file="dbdump/$MYSQL_SQL_FILENAME_$_now.sql"

echo $_file

#mysqldump "$MYSQL_DATABASE" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" > $_file