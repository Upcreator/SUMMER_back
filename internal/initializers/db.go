package initializers

import (
	"fmt"
	"log"

	"github.com/Upcreator/SUMMER_back/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection fail \n", err.Error())
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	DB.AutoMigrate(
		&models.NewsModel{},
		&models.TransitionApplicationModel{},
		&models.Election{},
		&models.Vote{},
		&models.Question{},
		&models.User{},
	)

	log.Println("Connected Successfully to DB")
}
