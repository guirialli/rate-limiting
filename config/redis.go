package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Redis struct {
	Addr     string `json:"REDIS_ADDR"`
	Password string `json:"REDIS_PASSWORD"`
	Db       int    `json:"REDIS_DB"`
}

func LoadRedisConfig() *Redis {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Db:       db,
	}
}
