package main

import (
	"archive/zip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/robfig/cron"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type uploadConfiguration struct {
	awsAccessKeyId     string
	awsSecretAccesKey  string
	awsDefaultRegion   string
	awsBucket          string
	minioUrl           string
	minioEnabled       bool
	minioSsl           bool
	targetFolderPrefix string
}

type mysqlBackupConfiguration struct {
	host     string
	port     string
	database string
	user     string
	password string
	allDbs   bool
}

type cleanConfiguration struct {
	cleanDays int
}

func main() {

	minioEnabledLocal, _ := strconv.ParseBool(os.Getenv("MINIO_ENABLED"))
	minioSslLocal, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))
	// Initialize variables
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

	// Initialize cron
	c := cron.New()
	_ = c.AddFunc(os.Getenv("SCHEDULE"), func() { fmt.Println("Running Scheduled Job: " + os.Getenv("SCHEDULE")) })
	_ = c.AddFunc(os.Getenv("SCHEDULE"), func() { runBackup(mysqlBackupConfiguration, uploadConfiguration) })
	c.Start()

	// Run main forever
	select {}

}

//https://github.com/minio/cookbook/blob/master/docs/aws-sdk-for-go-with-minio.md

func runBackup(mysqlBackupConfiguration mysqlBackupConfiguration, uploadConfiguration uploadConfiguration) {

	fmt.Println("Running backup operation...")
	uploadMinio(uploadConfiguration, backupMysql(mysqlBackupConfiguration))

}

func backupMysql(mysqlBackupConfiguration mysqlBackupConfiguration) string {

	//--skip_add_locks --skip-lock-tables --max_allowed_packet=1500M --complete-insert

	fmt.Printf("Executing mysqldump on %s database %s...\n", mysqlBackupConfiguration.host, mysqlBackupConfiguration.database)

	cmd := exec.Command("mysqldump",
		"--complete-insert",
		"--skip_add_locks",
		"--skip-lock-tables",
		"--max_allowed_packet=1500M",
		"-P"+mysqlBackupConfiguration.port,
		"-h"+mysqlBackupConfiguration.host,
		"-u"+mysqlBackupConfiguration.user,
		"-p"+mysqlBackupConfiguration.password,
		mysqlBackupConfiguration.database)

	//command :=  []string { "mysqldump",
	//	"-P"+mysqlBackupConfiguration.port,
	//	"-h"+mysqlBackupConfiguration.host,
	//	"-u"+mysqlBackupConfiguration.user,
	//	"-p"+mysqlBackupConfiguration.password,
	//	mysqlBackupConfiguration.database }
	//
	//fmt.Printf("command: %s \n",strings.Join(command," "))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now().Local()
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	fileDump := currentTime.Format("2006-01-02") + "-" + timestamp + "_dump.sql"

	err = ioutil.WriteFile(fileDump, bytes, 0644)
	if err != nil {
		panic(err)
	}

	//zip file before uploading
	filesToZip := []string{fileDump}
	fileDumpZip := currentTime.Format("2006-01-02") + "-" + timestamp + "_dump.zip"

	if err := zipFiles(fileDumpZip, filesToZip); err != nil {
		panic(err)
	}
	deleteFile(fileDump)

	return fileDumpZip

}

func uploadMinio(uploadConfiguration uploadConfiguration, filenameToUpload string) {

	file, err := os.Open(filenameToUpload)
	if err != nil {
		fmt.Println("Failed to open file", filenameToUpload, err)
		os.Exit(1)
	}
	defer file.Close()

	conf := aws.Config{
		Credentials:      credentials.NewStaticCredentials(uploadConfiguration.awsAccessKeyId, uploadConfiguration.awsSecretAccesKey, ""),
		Endpoint:         aws.String(uploadConfiguration.minioUrl),
		Region:           aws.String("eu-west-1"),
		DisableSSL:       aws.Bool(uploadConfiguration.minioSsl),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)

	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(uploadConfiguration.awsBucket),
		Key:    aws.String(uploadConfiguration.targetFolderPrefix + filenameToUpload),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	deleteFile(filenameToUpload)

	fmt.Printf("Successfully uploaded %s to %s\n", filenameToUpload, result.Location)

}

func deleteFile(path string) {
	err := os.Remove(path)

	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

func uploadS3(filename string, bucket string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	//Aws config from environment variables
	conf := aws.Config{}

	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)

	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)

}

func zipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
