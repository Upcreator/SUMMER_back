package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbUser     string
	dbPassword string
	dbName     string
	dbHost     string
	dbPort     string
)

func loadEnvVariables() error {
	godotenv.Load()

	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	return nil
}

var db *gorm.DB

func setupDB() {
	env_err := loadEnvVariables()
	if env_err != nil {
		log.Fatalf("Error loading env: %v", env_err)
	}

	// Create connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to DB using GORM
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Console log
	log.Println("DB connection: success")
}
