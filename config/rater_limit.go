package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type RaterLimit struct {
	IpRefresh    time.Duration
	JwtRefresh   time.Duration
	IpTrysMax    int
	JwtTrysMax   int
	BlockTimeout time.Duration
}

func LoadRaterLimitConfig() (*RaterLimit, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	jwtRefresh, err := strconv.Atoi(os.Getenv("JWT_REFRESH_ACCESS"))
	if err != nil {
		return nil, err
	}

	ipRefresh, err := strconv.Atoi(os.Getenv("IP_REFRESH_ACCESS"))
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
		JwtRefresh:   time.Duration(jwtRefresh) * time.Second,
		IpRefresh:    time.Duration(ipRefresh) * time.Second,
		JwtTrysMax:   jwtTryMax,
		BlockTimeout: time.Duration(blockTimeout) * time.Minute,
		IpTrysMax:    ipTrysMax,
	}, nil
}
