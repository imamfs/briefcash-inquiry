package config

import (
	"fmt"
	"os"

	logs "briefcash-inquiry/internal/helper/loghelper"

	env "github.com/joho/godotenv"
)

type Config struct {
	DBUrl        string
	DBHost       string
	DBUsername   string
	DBPassword   string
	DBPort       string
	DBName       string
	AppPort      string
	RedisAddress string
	RedisPort    string
}

func LoadConfig() (*Config, error) {
	if err := env.Load(); err != nil {
		logs.Logger.Error("No .env file found, using system environment variables")
	}

	cfg := &Config{
		DBUrl:        os.Getenv("DB_URL"),
		DBHost:       os.Getenv("DB_HOST"),
		DBUsername:   os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		RedisAddress: os.Getenv("REDIS_ADDRESS"),
		RedisPort:    os.Getenv("REDIS_PORT"),
		AppPort: func() string {
			if value := os.Getenv("APP_PORT"); value != "" {
				return value
			}
			return ":8080"
		}(),
	}

	if cfg.DBHost == "" {
		logs.Logger.Error("DB_HOST is not set in environment")
		return nil, fmt.Errorf("DB_HOST is not set in environment")
	}

	return cfg, nil
}
