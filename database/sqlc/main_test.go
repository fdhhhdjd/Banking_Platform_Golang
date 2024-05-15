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
	dbSource = "postgresql://taidev:eGMWFMa5gaAX1nfS229w@localhost:54333/service_banking?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	if err = testDB.Ping(); err != nil {
		log.Fatal("Can't ping the database: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
