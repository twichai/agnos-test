package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultDBURL = "postgres://postgres:postgres@localhost:5432/his?sslmode=disable"

func NewPostgresGormFromEnv() (*gorm.DB, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = defaultDBURL
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
