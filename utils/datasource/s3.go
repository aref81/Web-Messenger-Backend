package datasource

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func ConnectS3(accessKey, secretKey, region, endpoint string) (*session.Session, error) {
	s3Session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		log.Warnln("can not connect to s3", err)
		return nil, err
	}

	log.Infoln("connected to S3 instance")
	return s3Session, nil
}

func UploadS3(sess *session.Session, fileHeader *multipart.FileHeader, bucket string, ID string) (string, error) {
	uploader := s3manager.NewUploader(sess)
	file, err := fileHeader.Open()
	if err != nil {
		log.Warnln("cant open file")
		return "", err
	}

	fullFileName := ID + fileHeader.Filename
	key := fmt.Sprintf("%s", fullFileName)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Warnf("Unable to upload %q to %q, %v", fullFileName, bucket, err)
		return "", err
	}
	log.Infof("Successfully uploaded %q to %q\n", fullFileName, bucket)

	return key, nil
}

func DownloadS3(sess *session.Session, bucket string, key string) (*os.File, error) {
	getwd, _ := os.Getwd()

	file, err := os.Create(getwd + "/resources/profile/" + key)
	if err != nil {
		log.Warnf("Unable to open file %q, %v", key, err)
		return nil, err
	}
	defer file.Close()

	s3Client := s3.New(sess)

	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Warnf("Unable to download item %q, %v", key, err)
		return nil, err
	}

	_, err = io.Copy(file, obj.Body)
	if err != nil {
		log.Warnln("cant copy file")
		return nil, err
	}

	return file, nil
}
