package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type JWT struct {
	Secret   string
	ExpireIn int
	UnitTime rune
}

func LoadJwtConfig() *JWT {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	expireIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_IN"))
	if err != nil {
		panic(fmt.Errorf("JWT_EXPIRE_IN env variable not set int: %s", err.Error()))
	}

	return &JWT{
		Secret:   os.Getenv("JWT_SECRET"),
		ExpireIn: expireIn,
		UnitTime: rune(os.Getenv("JWT_UNIT_TIME")[0]),
	}
}
