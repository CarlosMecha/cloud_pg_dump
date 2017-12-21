package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/users"
)

type dropboxUploader struct {
	client files.Client
	users  users.Client
}

// NewDropboxUploader creates uploader connected to s3
func NewDropboxUploader(token string) Uploader {

	config := dropbox.Config{Token: token}

	return &dropboxUploader{
		client: files.New(config),
		users:  users.New(config),
	}
}

func (u *dropboxUploader) UploadFile(conf configuration, file string) (string, error) {

	if _, err := u.users.GetCurrentAccount(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Dropbox token: %s\n", err.Error())
		return "", err
	}

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	key := conf.dropboxPath + time.Now().Format("2006-01-02-15-04.backup")

	info := files.NewCommitInfo(key)
	info.Autorename = true

	if _, err := u.client.Upload(info, f); err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "File uploaded to dropbox:%s\n", key)
	return key, nil
}
