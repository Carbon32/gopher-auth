package db

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	JWTSecret string
}

var Environement = InitConfig()

func InitConfig() Config {
	// The default values are simply the values used for testing
	godotenv.Load()
	return Config{
		Host:      getEnv("PUBLIC_HOST", "localhost"),
		Port:      getNumericEnv("PORT", 5432),
		User:      getEnv("USER", "postgres"),
		Password:  getEnv("PASSWORD", "carbon123"),
		Database:  getEnv("DATABASE", "go_backend"),
		JWTSecret: getEnv("JWT_SECRET", "whateverman"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getNumericEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		number, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return number
	}

	return fallback
}
