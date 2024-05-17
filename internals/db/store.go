package db

import (
	"log"

	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
)

var storeInstance *database.Store

// InitStore initializes the database and store instance.
func InitStore() error {
	err := InitDB()
	if err != nil {
		return err
	}

	db := GetDB()
	storeInstance = database.NewStore(db)
	return nil
}

// GetStore returns the initialized store instance.
func GetStore() *database.Store {
	if storeInstance == nil {
		log.Fatal("Store is not initialized. Please call InitStore first.")
	}
	return storeInstance
}
