package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx"
)

const (
	TestHost     = "localhost"
	TestPort     = 5432
	TestUsername = "postgres"
	TestPassword = ""
	TestDatabase = "postgres"
)

func testConn(t *testing.T) *pgx.Conn {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     TestHost,
		Port:     TestPort,
		User:     TestUsername,
		Password: TestPassword,
		Database: TestDatabase,
	})
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func loadTestData(t *testing.T) {
	conn := testConn(t)
	defer conn.Close()

	conn.Exec("DROP SCHEMA test CASCADE")
	if _, err := conn.Exec("CREATE SCHEMA test"); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec("CREATE TABLE test.table (ID VARCHAR)"); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec("INSERT INTO test.table (ID) VALUES (1), (2), (3)"); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec("CREATE TABLE test.users (ID VARCHAR)"); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec("INSERT INTO test.users (ID) VALUES (1), (2), (3)"); err != nil {
		t.Fatal(err)
	}
}

func TestRunPgDump(t *testing.T) {

	loadTestData(t)

	file, err := RunPgDump(configuration{
		host:     TestHost,
		port:     fmt.Sprintf("%d", TestPort),
		database: TestDatabase,
		username: TestUsername,
		password: TestPassword,
		args:     []string{"-h", TestHost, "-U", TestUsername},
	})
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(file)

}
