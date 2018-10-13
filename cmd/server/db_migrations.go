package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/lib/pq"
)

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
