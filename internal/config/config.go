package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_Host     string
	DB_Port     string
	DB_User     string
	DB_Password string
	DB_Name     string
}

func LoadConfig() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	godotenv.Load(".env." + env)

	config := &Config{
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_User:     os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME"),
	}

	if config.DB_Host == "" {
		log.Println("WARNING: DB_HOST is empty! Check your .env file and path")
	}

	return config
}
