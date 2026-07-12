package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Docs     DocsConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type DocsConfig struct {
	Host   string
	Schema string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	AccessSecret   string
	UserAccessTTL  int // минуты
	AdminAccessTTL int // минуты (обычно меньше, чем у пользователя)
}

func Load() *Config {
	_ = godotenv.Load() // Игнорируем ошибку, если .env нет

	return &Config{
		Docs: DocsConfig{
			Host:   os.Getenv("SWAGGER_HOST"),
			Schema: os.Getenv("SWAGGER_SCHEMES"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "agro"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:   getEnv("JWT_ACCESS_SECRET", "change-me-in-prod"),
			UserAccessTTL:  getEnvInt("JWT_USER_TTL_MINUTES", 1440), // 24 часа
			AdminAccessTTL: getEnvInt("JWT_ADMIN_TTL_MINUTES", 60),  // 1 час
		},
	}
}

func (d DatabaseConfig) DSN() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.DBName + "?sslmode=" + d.SSLMode
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
