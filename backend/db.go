package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:password@db:5432/quizdb?sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	// Retry ping with exponential backoff to allow Postgres container to finish initialization.
	maxAttempts := 10
	baseDelay := 500 * time.Millisecond
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			lastErr = nil
			break
		}
		lastErr = err
		sleep := baseDelay * (1 << uint(attempt-1))
		if sleep > 5*time.Second {
			sleep = 5 * time.Second
		}
		log.Printf("db ping attempt %d/%d failed: %v; retrying in %s", attempt, maxAttempts, err, sleep)
		time.Sleep(sleep)
	}
	if lastErr != nil {
		log.Fatalf("failed to ping db after retries: %v", lastErr)
	}
	// simple schema: store quiz JSON in jsonb column
	// ensure schema; use background context since ping is already successful
	_, err = db.ExecContext(context.Background(), `
	CREATE TABLE IF NOT EXISTS quizzes (
		id TEXT PRIMARY KEY,
		title TEXT,
		data JSONB,
		created_at TIMESTAMPTZ
	)`)
	if err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}
}
