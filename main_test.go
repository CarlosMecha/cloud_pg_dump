package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {

	cases := []struct {
		args                  []string
		password              string
		expectedError         string
		expectedConfiguration configuration
	}{
		// OK
		{
			args:     []string{"-a", "-h", "host", "-U", "user", "-p", "1234", "-d", "database", "--s3bucket=bucket", "--s3prefix=/path/"},
			password: "password",
			expectedConfiguration: configuration{
				host:     "host",
				port:     "1234",
				username: "user",
				password: "password",
				database: "database",
				s3bucket: "bucket",
				s3prefix: "/path/",
				args:     []string{"-a", "-h", "host", "-U", "user", "-p", "1234", "-d", "database"},
			},
		},
		// Defaults
		{
			args: []string{"-a", "--s3bucket=bucket"},
			expectedConfiguration: configuration{
				host:        "localhost",
				port:        "5432",
				username:    "postgres",
				password:    "",
				database:    "postgres",
				s3bucket:    "bucket",
				s3prefix:    "",
				dropboxPath: "",
				args:        []string{"-a"},
			},
		},
		// OK (long args)
		{
			args:     []string{"-a", "--host=host", "--username=user", "--port=1234", "--dbname=database", "--s3bucket=bucket", "--s3prefix=/path/", "--dropboxpath=/app/"},
			password: "password",
			expectedConfiguration: configuration{
				host:        "host",
				port:        "1234",
				username:    "user",
				password:    "password",
				database:    "database",
				s3bucket:    "bucket",
				s3prefix:    "/path/",
				dropboxPath: "/app/",
				args:        []string{"-a", "--host=host", "--username=user", "--port=1234", "--dbname=database"},
			},
		},
		// file provided
		{
			args:          []string{"-a", "-f", "file", "--host=host", "--username=user", "--port=1234", "--dbname=database", "--s3bucket=bucket", "--s3prefix=/path/"},
			expectedError: "the argument --file or -f is not supported",
		},
		// -W provided
		{
			args:          []string{"-a", "-W", "--host=host", "--username=user", "--port=1234", "--dbname=database", "--s3bucket=bucket", "--s3prefix=/path/"},
			expectedError: "password argument -W or --password not supported",
		},
		// Missing bucket
		{
			args:          []string{"--host=host", "--username=user", "--port=1234", "--dbname=database", "--s3prefix=/path/"},
			expectedError: "missing s3 bucket",
		},
	}

	for i, test := range cases {
		os.Setenv("PGPASSWORD", test.password)
		conf, err := parseArguments(test.args)
		if err != nil {
			if test.expectedError != err.Error() {
				t.Errorf("Test %d: Expected error %s, got %s", i, test.expectedError, err.Error())
			}
			continue
		} else if test.expectedError != "" {
			t.Errorf("Test %d: Expected error %s, got nothing", i, test.expectedError)
		}

		assert.EqualValues(t, test.expectedConfiguration, conf, "Test %d: Expected configuration %+v, got %+v", i, test.expectedConfiguration, conf)
	}
}

func TestMain_integration(t *testing.T) {
	if os.Getenv("DROPBOX_TOKEN") == "" || os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("S3_TEST_BUCKET") == "" || os.Getenv("S3_TEST_REGION") == "" {
		t.Skip("No credentials provided for integration tests, skipping...")
	}

	os.Setenv("AWS_REGION", os.Getenv("S3_TEST_REGION"))

	loadTestData(t)

	command := fmt.Sprintf(
		"cloud_pg_dump --host=%s --username=%s --port=%d --dbname=%s --s3bucket=%s --s3prefix=%s --dropboxpath=%s",
		TestHost, TestUsername, TestPort, TestDatabase, os.Getenv("S3_TEST_BUCKET"), TestS3Path, TestDropboxPath,
	)

	os.Args = strings.Split(command, " ")

	main()
}

func TestMain(t *testing.T) {

	u := &fakeUploader{}
	Uploaders = []Uploader{u}

	loadTestData(t)

	command := fmt.Sprintf(
		"cloud_pg_dump --host=%s --username=%s --port=%d --dbname=%s --s3bucket=%s --s3prefix=%s --dropboxpath=%s",
		TestHost, TestUsername, TestPort, TestDatabase, os.Getenv("S3_TEST_BUCKET"), TestS3Path, TestDropboxPath,
	)

	os.Args = strings.Split(command, " ")

	main()

	if len(u.calls) != 1 {
		t.Fatal("Expected call to the uploader")
	}
}
