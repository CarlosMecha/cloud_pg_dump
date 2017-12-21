package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

// RunPgDump runs pg_dump as a subprocess and returns the file where the output is stored
func RunPgDump(conf configuration) (string, error) {

	if conf.password != "" {
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		pgpass, err := os.OpenFile(u.HomeDir+"/.pgpass", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return "", err
		}

		content := fmt.Sprintf("%s:%s:%s:%s:%s\n", conf.host, conf.port, conf.database, conf.username, conf.password)
		if _, err := pgpass.Write([]byte(content)); err != nil {
			return "", err
		}
		if err := pgpass.Close(); err != nil {
			return "", err
		}
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	file := f.Name()

	args := make([]string, 0, len(conf.args)+2)
	args = append(args, "-f", file)
	cmd := exec.Command("pg_dump", append(args, conf.args...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		os.Remove(file)
		return "", err
	}

	return file, nil
}
