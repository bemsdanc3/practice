package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"tickets/internal/entities"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&entities.Segment{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connection established and migration completed")
}

func GetDB() *gorm.DB {
	return DB
}
