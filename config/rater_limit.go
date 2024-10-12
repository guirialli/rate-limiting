package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type RaterLimit struct {
	IpTimeout    time.Duration `json:"IP_TIMEOUT"`
	JwtTimeout   time.Duration `json:"JWT_TIMEOUT"`
	IpTrysMax    int           `json:"IP_TRYS_MAX"`
	JwtTrysMax   int           `json:"JWT_TRYS_MAX"`
	BlockTimeout time.Duration `json:"BLOCK_TIMEOUT"`
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

	jwtTryMax, err := strconv.Atoi(os.Getenv("JWT_TRYS_MAX"))
	if err != nil {
		return nil, err
	}

	ipTrysMax, err := strconv.Atoi(os.Getenv("IP_TRYS_MAX"))
	if err != nil {
		return nil, err
	}

	blockTimeout, err := strconv.Atoi(os.Getenv("BLOCK_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	return &RaterLimit{
		JwtTimeout:   time.Duration(jwtTimeout) * time.Second,
		IpTimeout:    time.Duration(ipTimeout) * time.Second,
		JwtTrysMax:   jwtTryMax,
		BlockTimeout: time.Duration(blockTimeout) * time.Minute,
		IpTrysMax:    ipTrysMax,
	}, nil
}
