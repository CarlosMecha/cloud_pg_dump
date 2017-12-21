package main

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

type fakeUploader struct {
	calls []struct {
		config configuration
		file   string
	}
}

func (u *fakeUploader) UploadFile(conf configuration, file string) (string, error) {

	u.calls = append(u.calls, struct {
		config configuration
		file   string
	}{conf, file})

	return file, nil
}

func createTestFile(t *testing.T) string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s: This is a test", time.Now())); err != nil {
		t.Fatal(err)
	}

	return file.Name()
}
