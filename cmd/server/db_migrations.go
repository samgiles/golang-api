package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

func MigrateDatabaseUp(dbname string, db *sql.DB) error {

    assetSource := bindata.Resource(AssetNames(), Asset)

    sourceDriver, err :=  bindata.WithInstance(assetSource)

    if err != nil {
        return err
    }

    dbDriver, err := postgres.WithInstance(db, &postgres.Config{})

    if err != nil {
        return err
    }

    migrate, err := migrate.NewWithInstance("go-bindata", sourceDriver, dbname, dbDriver)

    if err != nil {
        return err
    }

    return migrate.Up()
}
