package sqlite3

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	// Imported to execute the inner init() function
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Imported to execute the inner init() function
	_ "github.com/mattn/go-sqlite3"
)

const schemaVersion = 1
const databaseName = ""
const driverName = "sqlite3"

// Config ...
type Config struct {
	DatabaseFilePath    string
	MigrationsDirectory string
}

// Database ...
type Database struct {
	Connection          *sql.DB
	MigrationsDirectory string
}

// NewDatabase opens a SQL connection pool and returns a Database instance
func NewDatabase(conf Config) (Database, error) {
	databaseFilePath, _ := filepath.Abs(conf.DatabaseFilePath)

	db, err := sql.Open(driverName, databaseFilePath)
	return Database{
		Connection:          db,
		MigrationsDirectory: conf.MigrationsDirectory,
	}, err
}

// RunMigrations trys to run the .sql migrations contained within the configured db.MigrationsDirectory
func (db Database) RunMigrations() error {
	driver, err := sqlite3.WithInstance(db.Connection, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance, %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", db.MigrationsDirectory),
		databaseName, driver)
	if err != nil {
		return fmt.Errorf("error running the database migrations, %w", err)
	}

	return m.Up()
}
