package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type s3Uploader struct {
	client s3iface.S3API
}

// NewS3Uploader creates uploader connected to s3
func NewS3Uploader(region string) Uploader {

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	return &s3Uploader{
		client: s3.New(sess),
	}
}

func (u *s3Uploader) UploadFile(conf configuration, file string) (string, error) {

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	key := conf.s3prefix + time.Now().Format("2006-01-02-15-04.backup")

	if _, err := u.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.s3bucket),
		Key:    aws.String(key),
		Body:   f,
	}); err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "File uploaded to s3://%s/%s\n", conf.s3bucket, key)
	return key, nil
}
