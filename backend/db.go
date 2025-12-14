package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:password@db:5432/quizdb?sslmode=disable"
	}

	var err error
	// open with retries similar to prior ping loop
	maxAttempts := 10
	baseDelay := 500 * time.Millisecond
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		sleep := baseDelay * (1 << uint(attempt-1))
		if sleep > 5*time.Second {
			sleep = 5 * time.Second
		}
		log.Printf("db open attempt %d/%d failed: %v; retrying in %s", attempt, maxAttempts, err, sleep)
		time.Sleep(sleep)
	}
	if err != nil {
		log.Fatalf("failed to open db after retries: %v", err)
	}

	// Auto-migrate the schema for Quiz and related tables (including new auth/submission models)
	if err := db.AutoMigrate(&Quiz{}, &Question{}, &Option{}, &Response{}, &User{}, &Submission{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
