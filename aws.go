package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadfile() {
	var AnonymousCredentials = credentials.NewStaticCredentials(idiam, secretiam, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: AnonymousCredentials,
		Region:      aws.String(regionaws),
	}))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	filename := localfile
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(bucketfile),
		Body:   f,
	})
	if err != nil {
		log.Printf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
}

func uploadvar() {
	var AnonymousCredentials = credentials.NewStaticCredentials(idiam, secretiam, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: AnonymousCredentials,
		Region:      aws.String(regionaws),
	}))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	filename := bucketfile
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(bucketfile),
		Body:   f,
	})
	if err != nil {
		log.Printf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
}

func uploadkey() {
	/*
		TODO:
	*/
}
