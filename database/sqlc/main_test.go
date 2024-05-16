package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"

	config "github.com/fdhhhdjd/Banking_Platform_Golang/configs"
	util "github.com/fdhhhdjd/Banking_Platform_Golang/internals/utils"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config.LoadConfig("../../configs")

	dbSource := util.P("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.AppConfig.Database.User, config.AppConfig.Database.Password, config.AppConfig.Database.Host, strconv.Itoa(config.AppConfig.Database.Port), config.AppConfig.Database.Name, config.AppConfig.Database.Ssl)

	log.Println(dbSource)

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
