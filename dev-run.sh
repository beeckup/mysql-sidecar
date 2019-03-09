#!/bin/sh

SCHEDULE="*/10 * * * * *" \
  AWS_ACCESS_KEY_ID="123123123" \
  AWS_SECRET_ACCESS_KEY="456456456456" \
  AWS_DEFAULT_REGION="eu-west-1" \
  DEV_FILENAME_TO_UPLOAD="testfile.txt" \
  AWS_S3_TARGET_BUCKET="testbucket" \
  AWS_S3_TARGET_BUCKET="testbucket" \
  MINIO_ENABLED="true" \
  MINIO_URL="http://localhost:9000" \
  MINIO_SSL="false" \
  TARGET_FOLDER_PREFIX="dump_database/mysqlbackup_" \
  MYSQL_HOST="127.0.0.1" \
  MYSQL_PORT="3306" \
  MYSQL_DATABASE="wordpress" \
  MYSQL_USER="root" \
  MYSQL_PASSWORD="123456" \
  MYSQL_ALL_DB="false" \
  go run backup.go