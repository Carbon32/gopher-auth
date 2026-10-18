package db

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string

	JWTSecret string

	SMTPHost           string
	SMTPPort           int
	SMTPUsername       string
	SMTPPassword       string
	SMTPUseTLS         bool
	SMTPFrom           string
	PublicAPIURL       string
	PublicAppURL       string
	HTTPAddress        string
	TLSCertFile        string
	TLSKeyFile         string
	CORSAllowedOrigins []string
}

var Environement = InitConfig()

func InitConfig() Config {
	// The default values are simply the values used for testing
	godotenv.Load()
	return Config{
		Host:     getEnv("DB_HOST", getEnv("PUBLIC_HOST", "localhost")),
		Port:     getNumericEnv("DB_PORT", getNumericEnv("PORT", 5432)),
		User:     getEnv("DB_USER", getEnv("USER", "postgres")),
		Password: getEnv("DB_PASSWORD", getEnv("PASSWORD", "carbon123")),
		Database: getEnv("DB_NAME", getEnv("DATABASE", "go_backend")),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "whateverman"),

		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getNumericEnv("SMTP_PORT", 1025),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPUseTLS:   getBoolEnv("SMTP_TLS", false),

		SMTPFrom: getEnv("SMTP_FROM", "noreply@localhost"),

		PublicAPIURL: getEnv("PUBLIC_API_URL", "http://localhost:8080/api/v1"), // Will be used for e-mail stuff
		PublicAppURL: getEnv("PUBLIC_APP_URL", "http://localhost:5173"),        // This will be used for password stuff

		HTTPAddress: getEnv("HTTP_ADDR", ":8080"),
		TLSCertFile: getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:  getEnv("TLS_KEY_FILE", ""),

		CORSAllowedOrigins: splitList(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")), // Cleaner design
	}
}

func splitList(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			result = append(result, value)
		}
	}
	return result
}

func getBoolEnv(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
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
