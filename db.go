package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Set the slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound errors for logger
			Colorful:                  true,                   // Disable color
		},
	)

	// Create connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to DB using GORM
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Console log
	log.Println("DB connection: success")

	// Call the auto-migrate function here
	autoMigrateDB()
}

func autoMigrateDB() {
	err := db.AutoMigrate(
		&models.User{},
		&models.Election{},
		&models.Question{},
		&models.Vote{},
		&models.TransitionApplication{},
		&models.News{},
	)
	if err != nil {
		log.Fatalf("Error during auto migration: %v", err)
	}

	log.Println("Database auto migration: success")
}
