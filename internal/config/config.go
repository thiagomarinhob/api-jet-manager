// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Configurações do servidor
	ServerAddress string
	GinMode       string

	// Configurações do banco de dados
	BLUEPRINT_DB_HOST     string
	BLUEPRINT_DB_PORT     int
	BLUEPRINT_DB_USERNAME string
	BLUEPRINT_DB_PASSWORD string
	BLUEPRINT_DB_DATABASE string
	DBSSLMode             string

	// Configurações do JWT
	JWTSecret     string
	JWTExpiration time.Duration
}

func LoadConfig() (*Config, error) {
	// Carrega variáveis de ambiente do arquivo .env se existir
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))

	return &Config{
		// Servidor
		ServerAddress: getEnv("SERVER_ADDRESS", "8080"),
		GinMode:       getEnv("GIN_MODE", "debug"),

		// Banco de dados
		BLUEPRINT_DB_HOST:     getEnv("DB_HOST", "psql_bp"),
		BLUEPRINT_DB_PORT:     dbPort,
		BLUEPRINT_DB_USERNAME: getEnv("DB_USER", "docker"),
		BLUEPRINT_DB_PASSWORD: getEnv("DB_PASSWORD", "docker"),
		BLUEPRINT_DB_DATABASE: getEnv("DB_NAME", "jetmanager"),
		DBSSLMode:             getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "25thiago99"),
		JWTExpiration: time.Duration(jwtExpiration) * time.Hour,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
