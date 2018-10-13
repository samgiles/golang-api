package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func CreateDatabaseConnection(user, password, host, dbname string) (*sql.DB, error) {
	return sql.Open("postgres", createConnectionString(user, password, host, dbname))
}

func WaitForDbConnectivity(db *sql.DB, timeout time.Duration) error {
	connectedAck := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			err := db.Ping()
			if err == nil {
				connectedAck <- true
				return
			}

			log.Printf("DB connection err: %s", err.Error())
		}
	}()

	select {
	case <-connectedAck:
		return nil
	case <-time.After(timeout):
		return NewOperationTimeoutError("Timed out waiting for DB connection confirmation")
	}
}

func createConnectionString(user, password, host, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
}

func MigrateDatabaseUp(dbname string, db *sql.DB) error {

	assetSource := bindata.Resource(AssetNames(), Asset)

	sourceDriver, err := bindata.WithInstance(assetSource)

	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceDriver, dbname, dbDriver)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != migrate.ErrNoChange {
		return err
	}

	return nil
}
