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


if [ "$S3_UPLOAD" = "true" ]; then
echo "Devo caricarlo su minio..."

bucket="$S3_BUCKET"

host="$S3_HOST"
link="$S3_PROTOCOL""://""$S3_HOST"

echo "$link"

s3_key="$S3_KEY"
s3_secret="$S3_SECRET"

resource="/${bucket}/${_file}"
content_type="application/octet-stream"
date=`date -R`
_signature="PUT\n\n${content_type}\n${date}\n${resource}"
signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`
echo "Eseguo curl"
curl -v -X PUT -T "${_file}" \
          -H "Host: $host" \
          -H "Date: ${date}" \
          -H "Content-Type: ${content_type}" \
          -H "Authorization: AWS ${s3_key}:${signature}" \
           $link${resource}

rm $_file


fi