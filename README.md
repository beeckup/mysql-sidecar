<img src="./docimages/LOGO_oriz.png" alt="logo" width="350"/> <img src="./docimages/logo-mysql.png" alt="logo" width="350"/>

# Sidecar Mysql

## Automatic Mysql Backup and upload on AWS S3 or Minio

Goland container to schedule and backup mysql database

Example deploy on  ```deploy_sidecar_example/docker-compose.yml```

Copy `env.sample` as `.env`

ENVIROMENT VARIABLE   | DESCRIPTION | Values
----------   | ---------- | --------------  
AWS_ACCESS_KEY_ID | Aws access key or Minio username | Access key string
AWS_SECRET_ACCESS_KEY | Aws secret key or Minio password | Secret access key string
AWS_DEFAULT_REGION | Aws default region or any value for minio | `eu-west-1` etc
AWS_S3_TARGET_BUCKET | Aws bucket name or minio bucket name | `bucketname`
MINIO_ENABLED | If target upload is a minio server | `true` or `false`
MINIO_SSL | If minio is SSL | `true` or `false`
MINIO_URL | Minio url | `http://localhost:9000` like
MYSQL_HOST | hostname or ip server mysql | hostname or ip
MYSQL_PORT | mysql port | `3306` or custom port
MYSQL_DATABASE | database name | string
MYSQL_USER | database user | string
MYSQL_PASSWORD | database password | string
MYSQL_ALL_DB | cycle all database and backups single file each | `true` or `false`
SCHEDULE | see below | 
CLEAN_DAYS | number of backup retention days | integer or empty


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



# Usage

Create `.env` file:

```bash

```

Create `docker-compose.yml` file:

```yml

```

Launch with

```bash
docker-compose up -d
```
