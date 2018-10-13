package main

import (
	"database/sql"
	"fmt"
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
