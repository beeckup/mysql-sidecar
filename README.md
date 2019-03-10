<img src="./images/LOGO_oriz.png" alt="logo" height="150"/> <img src="./images/logo-mysql.png" alt="logo" height="150"/>

# Sidecar Mysql

## Automatic Mysql Backup and upload on AWS S3 or Minio

Goland container to schedule and backup mysql database

Examples deploy on  ```examples/```

Copy `env.sample` as `.env`

ENVIROMENT VARIABLE   | DESCRIPTION | Values
----------   | ---------- | --------------  
TARGET_FOLDER_PREFIX | folder and prefix filename pattern on upload S3 | `dump_database/prefix_` 
SCHEDULE | see below | `0 * * * * *` once per minute
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
CLEAN_DAYS | number of backup retention days | integer or empty
CLEAN_FOLDER | folder prefix to clean | eg `dump_database/`


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



# Usage, AWS Example

Go to `examples/aws-s3/`

Create file `.env` and configure parameters

Launch with

```bash
docker-compose up -d
```


# Usage, MINIO Example

Go to `examples/minio/`

Create file `.env` and configure parameters

Launch with

```bash
docker-compose up -d
```

# Kubernetes Chart repository

Go to [beeckup/charts](https://github.com/beeckup/charts) for kubernetes chart