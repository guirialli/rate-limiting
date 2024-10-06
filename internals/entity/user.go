package entity

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func isValidPassword(password string) bool {
	var low, upper, digit, special bool
	if len(password) < 8 {
		return false
	}

	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			digit = true
		case unicode.IsLower(r):
			low = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			special = true
		case low && upper && digit && special:
			break
		}

	}

	return low && upper && digit && special
}

func NewUser(username, password string) (*User, error) {
	if !isValidPassword(password) {
		return nil, fmt.Errorf("password is very weak")
	} else if len(username) > 64 || len(username) <= 3 {
		return nil, fmt.Errorf("invalid username %s", username)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	return &User{
		Id:       uuid.NewString(),
		Username: username,
		Password: string(passwordHash),
	}, nil
}
