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

func uploadS3(uploadConfiguration uploadConfiguration, filenameToUpload filename) {

	fmt.Println("Uploading file on S3...")

	file, err := os.Open(string(filenameToUpload))
	if err != nil {
		fmt.Println("Failed to open file", filenameToUpload, err)
		os.Exit(1)
	}
	defer file.Close()

	conf := aws.Config{
		Credentials: credentials.NewStaticCredentials(uploadConfiguration.awsAccessKeyId, uploadConfiguration.awsSecretAccesKey, ""),
		Endpoint:    aws.String(uploadConfiguration.minioUrl),
		Region:      aws.String(uploadConfiguration.awsDefaultRegion),
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
