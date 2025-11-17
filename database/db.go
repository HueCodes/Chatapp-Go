package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	log.Printf("Database initialized: %s", dbPath)
	return nil
}

// AutoMigrate runs auto migration for the given models
func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}
