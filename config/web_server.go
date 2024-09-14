package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type WebServer struct {
	Ip   string
	Port int
}

func LoadWebServerConfig() (*WebServer, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	ip := os.Getenv("IP")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	return &WebServer{
		Ip:   ip,
		Port: port,
	}, nil
}
