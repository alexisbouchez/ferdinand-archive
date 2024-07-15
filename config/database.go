package config

import (
	"errors"
	"ferdinand/app/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDatabase opens a connection to the database and returns a pointer to the gorm.DB object.
func NewDatabase() (*gorm.DB, error) {
	// Get the appropriate gorm.Dialector based on the DSN environment variable
	dialector := getDialector()
	if dialector == nil {
		return nil, errors.New("invalid DBMS")
	}

	// Open a database connection
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.AutoMigrate(
		&models.User{},
		&models.Domain{},
		&models.APIKey{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

// getDialector returns the appropriate gorm.Dialector (SQLite, PostgreSQL, or MySQL)
// based on the DBMS environment variable.
func getDialector() gorm.Dialector {
	switch os.Getenv("DBMS") {
	case "sqlite":
		return sqlite.Open(os.Getenv("DSN"))
	case "postgres":
		return postgres.Open(os.Getenv("DSN"))
	case "mysql":
		return mysql.Open(os.Getenv("DSN"))
	default:
		return nil
	}
}
