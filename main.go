package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const usage = `
s3pgdump v1.0.0

env:
  - PGPASSWORD=<password>
  - DROPBOX_TOKEN=<token>
  - AWS_REGION=<region>
  - AWS_<creds>
USAGE: s3pgdump <pg_dump args> --s3bucket=bucket --s3prefix=/prefix --dropboxpath=/path
`

type configuration struct {
	host        string
	port        string
	database    string
	username    string
	password    string
	s3bucket    string
	s3prefix    string
	dropboxPath string
	args        []string
}

// parseArguments returns the configuration, including pg_dump arguments
func parseArguments(args []string) (configuration, error) {

	conf := configuration{
		password: os.Getenv("PGPASSWORD"),
		args:     make([]string, 0),
	}

	for i, arg := range args {
		if arg == "-f" || strings.Contains(arg, "--file=") {
			return configuration{}, errors.New("the argument --file or -f is not supported")
		} else if arg == "-W" || arg == "--password" {
			return configuration{}, errors.New("password argument -W or --password not supported")
		} else if arg == "-h" && i < len(args)-1 {
			conf.host = args[i+1]
		} else if strings.Contains(arg, "--host=") {
			conf.host = strings.Split(arg, "=")[1]
		} else if arg == "-p" && i < len(args)-1 {
			conf.port = args[i+1]
		} else if strings.Contains(arg, "--port=") {
			conf.port = strings.Split(arg, "=")[1]
		} else if arg == "-d" && i < len(args)-1 {
			conf.database = args[i+1]
		} else if strings.Contains(arg, "--dbname=") {
			conf.database = strings.Split(arg, "=")[1]
		} else if arg == "-U" && i < len(args)-1 {
			conf.username = args[i+1]
		} else if strings.Contains(arg, "--username=") {
			conf.username = strings.Split(arg, "=")[1]
		} else if strings.Contains(arg, "--s3bucket=") {
			conf.s3bucket = strings.Split(arg, "=")[1]
			continue
		} else if strings.Contains(arg, "--s3prefix=") {
			conf.s3prefix = strings.Split(arg, "=")[1]
			continue
		} else if strings.Contains(arg, "--dropboxpath=") {
			conf.dropboxPath = strings.Split(arg, "=")[1]
			continue
		}
		conf.args = append(conf.args, arg)
	}

	if conf.host == "" {
		conf.host = "localhost"
	}
	if conf.port == "" {
		conf.port = "5432"
	}
	if conf.database == "" {
		conf.database = "postgres"
	}
	if conf.username == "" {
		conf.username = "postgres"
	}
	if conf.s3bucket == "" {
		return configuration{}, errors.New("missing s3 bucket")
	}
	return conf, nil

}

func main() {

	conf, err := parseArguments(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
		return
	}

	output, err := RunPgDump(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
		return
	}
	defer os.Remove(output)

	for _, u := range Uploaders {
		if _, err := u.UploadFile(conf, output); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(1)
		}
	}
}
