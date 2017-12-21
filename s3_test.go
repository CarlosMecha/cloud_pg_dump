package main

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const TestS3Path = "/cloudPgDump/"

func TestS3UploadFile(t *testing.T) {
	bucket := os.Getenv("S3_TEST_BUCKET")
	region := os.Getenv("S3_TEST_REGION")
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || bucket == "" || region == "" {
		t.Skip("No AWS credentials or bucket provided for S3 integration tests, skipping...")
	}

	file := createTestFile(t)
	defer os.Remove(file)

	uploader := NewS3Uploader(region).(*s3Uploader)
	conf := configuration{
		s3bucket: bucket,
		s3prefix: TestS3Path,
	}

	key, err := uploader.UploadFile(conf, file)
	if err != nil {
		t.Errorf("Error uploading file to S3: %s", err.Error())
	}

	if _, err := uploader.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		t.Errorf("Error uploading file to S3: %s", err.Error())
	}

}
