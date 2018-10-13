package main

import (
	"log"
	"os"
	"reflect"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	db, err := CreateDatabaseConnection(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbName)

	if err != nil {
		log.Fatalf("main: could not create db conn: %s", err)
	}

	log.Println("main: migrating database up..")
	err = MigrateDatabaseUp(dbName, db)

	if err != nil {
		log.Printf("main: type %s", reflect.TypeOf(err))
		log.Fatalf("main: could not migrate db: %s", err)
	}
	log.Println("main: migrated database up..")

	server := NewServer(db)
	defer server.Stop()
	server.Start()
}
