package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
	AppPort    string
}

// LoadConfig reads configuration from .env file and returns Config struct
func LoadConfig() *Config {
	// Load .env file
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	return &Config{
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBPort:  dbPort,
		DBUser:  getEnv("DB_USER", "root"),
		DBPass:  getEnv("DB_PASS", ""),
		DBName:  getEnv("DB_NAME", "article_db"),
		AppPort: getEnv("APP_PORT", "8080"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}