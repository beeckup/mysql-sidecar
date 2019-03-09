#!/bin/sh

SCHEDULE="*/10 * * * * *" \
  ONLY_ONCE="true" \
  AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
  AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
  AWS_S3_TARGET_BUCKET="nutellino.tests" \
  AWS_DEFAULT_REGION="eu-west-1" \
  MINIO_ENABLED="false" \
  TARGET_FOLDER_PREFIX="dump_database/mysqlbackup_" \
  MYSQL_HOST="127.0.0.1" \
  MYSQL_PORT="3306" \
  MYSQL_DATABASE="wordpress" \
  MYSQL_USER="root" \
  MYSQL_PASSWORD="123456" \
  MYSQL_ALL_DB="false" \
  go run backup.go