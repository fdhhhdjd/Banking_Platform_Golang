package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://taidev:eGMWFMa5gaAX1nfS229w@localhost:5432/service_banking?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
