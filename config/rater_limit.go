package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type RaterLimit struct {
	JwtTimeout int `json:"jwt_timeout"`
	IpTimeout  int `json:"ip_timeout"`
}

func LoadRaterLimitConfig() (*RaterLimit, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	jwtTimeout, err := strconv.Atoi(os.Getenv("JWT_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	ipTimeout, err := strconv.Atoi(os.Getenv("IP_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	return &RaterLimit{
		JwtTimeout: jwtTimeout,
		IpTimeout:  ipTimeout,
	}, nil
}
