// db.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbUser     string
	dbPassword string
	dbName     string
	dbHost     string
	dbPort     string
)
var db *sql.DB

// Load from .env or from env itself
func loadEnvVariables() error {
	godotenv.Load()

	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	return nil
}

func setupDB() *sql.DB {
	err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	// Create connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)

	// Connect to DB
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error DB connection: %v", err)
	}

	// Check the DB connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Console log
	log.Println("DB connection: success")
	return db
}
