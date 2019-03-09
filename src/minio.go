package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
)

func uploadMinio(uploadConfiguration uploadConfiguration, filenameToUpload filename) {

	fmt.Println("Uploading file on Minio...")

	file, err := os.Open(string(filenameToUpload))
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

	//fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(uploadConfiguration.awsBucket),
		Key:    aws.String(uploadConfiguration.targetFolderPrefix + filepath.Base(string(filenameToUpload))),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	filenameToUpload.delete()

	fmt.Printf("Successfully uploaded %s to %s\n", filenameToUpload, result.Location)

}
