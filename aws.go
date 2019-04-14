package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getCredentials() *session.Session {
	var AnonymousCredentials = credentials.NewStaticCredentials(idiam, secretiam, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: AnonymousCredentials,
		Region:      aws.String(regionaws),
	}))
	return sess
}

func uploadfile() {

	sess := getCredentials()
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	filename := localfile
	f, err := os.Open(filename)
	if err != nil {
		log.Printf(fileOpenError, filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(bucketfile),
		Body:   f,
	})
	if err != nil {
		log.Printf(fileUploadError, err)
	}
	fmt.Printf(fileUploadError2, result.Location)
}

func generalLambda(funcName string, funcParams string) string {
	resp := ""
	sess := getCredentials()
	svc := lambda.New(sess)
	//funcParams = "keyblc:testparamformGO"
	bytespayload, err := json.Marshal(funcParams)
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(funcName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        bytespayload,
	}

	result, err := svc.Invoke(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidRequestContentException:
				fmt.Println(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
			case lambda.ErrCodeRequestTooLargeException:
				fmt.Println(lambda.ErrCodeRequestTooLargeException, aerr.Error())
			case lambda.ErrCodeUnsupportedMediaTypeException:
				fmt.Println(lambda.ErrCodeUnsupportedMediaTypeException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
			case lambda.ErrCodeEC2UnexpectedException:
				fmt.Println(lambda.ErrCodeEC2UnexpectedException, aerr.Error())
			case lambda.ErrCodeSubnetIPAddressLimitReachedException:
				fmt.Println(lambda.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
			case lambda.ErrCodeENILimitReachedException:
				fmt.Println(lambda.ErrCodeENILimitReachedException, aerr.Error())
			case lambda.ErrCodeEC2ThrottledException:
				fmt.Println(lambda.ErrCodeEC2ThrottledException, aerr.Error())
			case lambda.ErrCodeEC2AccessDeniedException:
				fmt.Println(lambda.ErrCodeEC2AccessDeniedException, aerr.Error())
			case lambda.ErrCodeInvalidSubnetIDException:
				fmt.Println(lambda.ErrCodeInvalidSubnetIDException, aerr.Error())
			case lambda.ErrCodeInvalidSecurityGroupIDException:
				fmt.Println(lambda.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
			case lambda.ErrCodeInvalidZipFileException:
				fmt.Println(lambda.ErrCodeInvalidZipFileException, aerr.Error())
			case lambda.ErrCodeKMSDisabledException:
				fmt.Println(lambda.ErrCodeKMSDisabledException, aerr.Error())
			case lambda.ErrCodeKMSInvalidStateException:
				fmt.Println(lambda.ErrCodeKMSInvalidStateException, aerr.Error())
			case lambda.ErrCodeKMSAccessDeniedException:
				fmt.Println(lambda.ErrCodeKMSAccessDeniedException, aerr.Error())
			case lambda.ErrCodeKMSNotFoundException:
				fmt.Println(lambda.ErrCodeKMSNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidRuntimeException:
				fmt.Println(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return ""
	}
	s := string(result.Payload)
	mapBody := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(s), &mapBody)
	if err2 != nil {
		log.Fatal(awsError)
	}
	var statusCode float64
	body := ""
	statusCode = 0
	for _, value := range mapBody {
		switch v := value.(type) {
		case float64:
			statusCode = value.(float64)
		case string:
			body = value.(string)
		default:
			println(v)
			log.Fatal("Error in parser type resp aws")
		}
	}
	if statusCode != 0 && body != "" {
		if statusCode == 200 || statusCode == 201 || statusCode == 202 {
			resp = body
		} else {
			resp = "ko"
			fmt.Println(awsError)
		}
	}

	return resp
}
