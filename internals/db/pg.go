package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	config "github.com/fdhhhdjd/Banking_Platform_Golang/configs"
	util "github.com/fdhhhdjd/Banking_Platform_Golang/internals/utils"
	_ "github.com/lib/pq" // or another driver you're using
)

var db *sql.DB

const (
	dbDriver = "postgres"
)

func InitDB() error {
	var err error
	config.LoadConfig("../../configs")

	dataSourceName := util.P("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.AppConfig.Database.User, config.AppConfig.Database.Password, config.AppConfig.Database.Host, strconv.Itoa(config.AppConfig.Database.Port), config.AppConfig.Database.Name, config.AppConfig.Database.Ssl)

	db, err = sql.Open(dbDriver, dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("CONNECTED TO POSTGRESQL SUCCESS 🐘!!")
	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database is not initialized. Please call InitDB first.")
	}
	return db
}
