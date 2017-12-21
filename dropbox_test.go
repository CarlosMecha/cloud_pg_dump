package main

import (
	"os"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

const TestDropboxPath = "/cloudPgDump-"

func TestDropboxUpload(t *testing.T) {
	token := os.Getenv("DROPBOX_TOKEN")
	if token == "" {
		t.Skip("No Dropbox token provided for integration tests, skipping...")
	}

	file := createTestFile(t)
	defer os.Remove(file)

	uploader := NewDropboxUploader(token).(*dropboxUploader)
	conf := configuration{dropboxPath: TestDropboxPath}

	key, err := uploader.UploadFile(conf, file)
	if err != nil {
		t.Errorf("Error uploading file to Dropbox: %s", err.Error())
	}

	if _, err := uploader.client.GetMetadata(&files.GetMetadataArg{Path: key}); err != nil {
		t.Errorf("Error uploading file to Dropbox: %s", err.Error())
	}

}
