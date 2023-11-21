package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check and enable vector extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return nil, fmt.Errorf("failed to enable vector extension: %w", err)
	}

	// Verify vector extension is working
	var version string
	if err := db.Raw("SELECT extversion FROM pg_extension WHERE extname = 'vector'").Scan(&version).Error; err != nil {
		return nil, fmt.Errorf("vector extension not properly installed: %w", err)
	}
	log.Println("pgvector extension version: ", version)

	// Auto-migrate the schema
	if err := db.AutoMigrate(&Actor{}, &Session{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schemas: %w", err)
	}

	// Create fragment tables
	if err := CreateFragmentTables(db); err != nil {
		return nil, fmt.Errorf("failed to create fragment tables: %w", err)
	}

	return db, nil
}

func CreateFragmentTables(db *gorm.DB) error {
	for _, table := range fragmentTables {
		// Create table
		if !db.Migrator().HasTable(string(table)) {
			// Create table only if it doesn't exist
			if err := db.Table(string(table)).Migrator().CreateTable(&Fragment{}); err != nil {
				return fmt.Errorf("failed to create %s table: %w", table, err)
			}
		}
	}
	return nil
}
