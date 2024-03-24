package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

func InitializeORM() {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		logger.Fatalf("DB_DSN environment variable is not present.")
		return
	}

	instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gl.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gl.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gl.Warn,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  true,
			},
		),
	})

	db = instance

	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
		return
	}

	AutoMigrate()
}

func AutoMigrate() {
	db.AutoMigrate(&User{})
}
