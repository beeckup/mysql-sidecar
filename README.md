# Sidecar Backup Mysql

Example deploy on  ```deploy_sidecar_example/docker-compose.yml```

Copy `env.sample` as `.env`

ENVIROMENT VARIABLE   | DESCRIPTION | Values
----------   | ---------- | --------------  
MYSQL_HOST | hostname or ip server mysql | hostname or ip
MYSQL_DATABASE | database name | string
MYSQL_USER | database user | string
MYSQL_PASSWORD | database password | string
MYSQL_SQL_FILENAME | filename backup | string
SCHEDULE | see below | 
ZIP_FILE | true to enable tar.gz compression | `true` or empty

## Schedule

Tip on ```SCHEDULE``` enviroment variable:

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?



## Minio/S3 config

ENVIROMENT VARIABLE   | DESCRIPTION | Values
----------   | ---------- | --------------  
S3_UPLOAD | Flag to enable s3 upload | `true` or empty
S3_BUCKET | Bucket name | string
S3_HOST | host:port | `host:port`
S3_PROTOCOL | protocol type | `http` or `https`
S3_KEY | key | string
S3_SECRET | secret | string
MINIO_PORT | local minio port to expose on host | port number