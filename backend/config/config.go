package config

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

var AppConfig Config

func LoadConfig() {
	godotenv.Load()
	
	AppConfig = Config{
		Port:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://chess:chess@localhost:5432/chess?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open(AppConfig.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("Database connected successfully")
}