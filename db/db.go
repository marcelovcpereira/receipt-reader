package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	config := NewConfig()
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName)

	fmt.Printf("connecting to %s as %s\n", config.DbHost, config.DbUser)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = migrateDatabase(db.DB)
	if err != nil {
		return nil, err
	}

	return db, err
}

func migrateDatabase(db *sql.DB) error {
	path, err := migrationsPath()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Printf("successfully migrated database")
	return nil
}

func migrationsPath() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	fileDirectory := filepath.Dir(file)
	return "file:///" + fileDirectory + "/migrations", nil
}
