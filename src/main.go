package main

import (
	"fmt"
	"github.com/robfig/cron"
	"os"
	"strconv"
)

func main() {

	// Initialize variables

	minioEnabledLocal, _ := strconv.ParseBool(os.Getenv("MINIO_ENABLED"))
	minioSslLocal, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))

	uploadConfiguration := uploadConfiguration{
		awsAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		awsSecretAccesKey:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
		awsDefaultRegion:   os.Getenv("AWS_DEFAULT_REGION"),
		awsBucket:          os.Getenv("AWS_S3_TARGET_BUCKET"),
		minioUrl:           os.Getenv("MINIO_URL"),
		minioEnabled:       minioEnabledLocal,
		minioSsl:           minioSslLocal,
		targetFolderPrefix: os.Getenv("TARGET_FOLDER_PREFIX"),
	}

	mysqlAllDbsLocal, _ := strconv.ParseBool(os.Getenv("MYSQL_ALL_DB"))

	mysqlBackupConfiguration := mysqlBackupConfiguration{
		host:     os.Getenv("MYSQL_HOST"),
		port:     os.Getenv("MYSQL_PORT"),
		database: os.Getenv("MYSQL_DATABASE"),
		user:     os.Getenv("MYSQL_USER"),
		password: os.Getenv("MYSQL_PASSWORD"),
		allDbs:   mysqlAllDbsLocal,
	}

	// Clean configuration

	cleanDaysLocal, _ := strconv.ParseInt(os.Getenv("CLEAN_DAYS"), 10, 64)

	cleanConfiguration := cleanConfiguration{
		folder:    os.Getenv("CLEAN_FOLDER"),
		cleanDays: cleanDaysLocal,
	}

	onlyOnceLocal, _ := strconv.ParseBool(os.Getenv("ONLY_ONCE"))

	runConfiguration := runConfiguration{
		onlyOnce: onlyOnceLocal,
	}

	// Initialize cron
	c := cron.New()
	_ = c.AddFunc(os.Getenv("SCHEDULE"), func() { fmt.Println("Running Scheduled Job: " + os.Getenv("SCHEDULE")) })
	_ = c.AddFunc(os.Getenv("SCHEDULE"), func() { runBackup(mysqlBackupConfiguration, uploadConfiguration, runConfiguration, cleanConfiguration) })
	c.Start()

	// Run main forever

	select {}

}

func runBackup(mysqlBackupConfiguration mysqlBackupConfiguration, uploadConfiguration uploadConfiguration, runConfiguration runConfiguration, cleanConfiguration cleanConfiguration) {

	fmt.Println("Running backup operation...")

	if uploadConfiguration.minioEnabled {
		uploadMinio(uploadConfiguration, backupMysql(mysqlBackupConfiguration))
		cleanMinio(uploadConfiguration, cleanConfiguration)

	} else {
		uploadS3(uploadConfiguration, backupMysql(mysqlBackupConfiguration))
		cleanS3(uploadConfiguration, cleanConfiguration)

	}

	if runConfiguration.onlyOnce {
		os.Exit(0)
	}

}
