package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Database struct {
	Hostname string
	Port     string
	User     string
	Password string
	Database string
}

func LoadDatabaseConfig() (*Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Database{
		Hostname: os.Getenv("DB_HOSTNAME"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}, nil
}
