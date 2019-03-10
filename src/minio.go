package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
	"time"
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

	sess := session.Must(session.NewSession(&conf))
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

func cleanMinio(uploadConfiguration uploadConfiguration, cleanConfiguration cleanConfiguration) {

	if cleanConfiguration.cleanDays == 0 {
		fmt.Println("Cannot proceed with cleaning operation, folder pattern CLEAN_DAYS unset or zero!")
		return
	}

	if cleanConfiguration.folder == "" {
		fmt.Println("Cannot proceed with cleaning operation, folder pattern CLEAN_FOLDER empty!")
		return
	}

	fmt.Printf("Proceeding with cleaning operations folder pattern %s older than %d days...\n", cleanConfiguration.folder, cleanConfiguration.cleanDays)

	conf := aws.Config{
		Credentials:      credentials.NewStaticCredentials(uploadConfiguration.awsAccessKeyId, uploadConfiguration.awsSecretAccesKey, ""),
		Endpoint:         aws.String(uploadConfiguration.minioUrl),
		Region:           aws.String("eu-west-1"),
		DisableSSL:       aws.Bool(uploadConfiguration.minioSsl),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.Must(session.NewSession(&conf))
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(uploadConfiguration.awsBucket),
		Prefix: aws.String(cleanConfiguration.folder),
	}

	nowTime := time.Now()
	resp, _ := svc.ListObjects(params)
	for _, key := range resp.Contents {
		//fmt.Println(*key.Key)
		//fmt.Println(*key.LastModified)
		//fmt.Println(getDifferenceDays(&nowTime,key.LastModified))
		//fmt.Println(getDifferenceMinutes(&nowTime,key.LastModified))

		if getDifferenceDays(&nowTime, key.LastModified) > cleanConfiguration.cleanDays {
			_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(uploadConfiguration.awsBucket), Key: aws.String(*key.Key)})
			if err != nil {
				fmt.Println("error", err)
				os.Exit(1)
			}
			err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
				Bucket: aws.String(uploadConfiguration.awsBucket),
				Key:    aws.String(*key.Key),
			})
			fmt.Printf("Object %s deleted.\n", *key.Key)
		}

	}

	fmt.Println("Bucket/folder cleaned.")

}
