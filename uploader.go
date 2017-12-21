package main

import (
	"os"
)

// Uploader uploads the file to the cloud
type Uploader interface {

	// UploadFile reads and uploads the file. Returns the file name or key
	// if the file was succesfully uploader, or an error otherwise.
	UploadFile(configuration, string) (string, error)
}

// Uploaders is the list of cloud uploaders
var Uploaders = []Uploader{
	NewS3Uploader(os.Getenv("AWS_REGION")),
	NewDropboxUploader(os.Getenv("DROPBOX_TOKEN")),
}
