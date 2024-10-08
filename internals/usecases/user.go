package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/pkg/uow"
	"time"
	"unicode"
)

type User struct {
	secret   string
	expireIn int
	unitTime rune // s to seconds, m to minutes, h to hours or d to days
}

func NewUser(secret string, expireIn int, unitTime rune) (*User, error) {
	unitTime = unicode.ToLower(unitTime)
	if unitTime != 's' && unitTime != 'm' && unitTime != 'h' && unitTime != 'd' {
		return nil, fmt.Errorf("invalid unit time")
	}
	return &User{
		secret:   secret,
		expireIn: expireIn,
		unitTime: unitTime,
	}, nil
}

func (u *User) getExpirationTime() int64 {
	exp := time.Duration(u.expireIn)
	switch u.unitTime {
	case 'd':
		exp = exp * 24 * time.Hour
	case 'h':
		exp = exp * time.Hour
	case 'm':
		exp = exp * time.Minute
	case 's':
		exp = exp * time.Second
	}

	return time.Now().Add(exp).Unix()
}

func (u *User) NewTokenAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(u.secret), nil)
}
func (u *User) genJwt(user *entity.User) (string, error) {
	token := u.NewTokenAuth()
	_, tokenString, err := token.Encode(map[string]interface{}{
		"sub":      user.Id,
		"username": user.Username,
		"exp":      u.getExpirationTime(),
	})
	if err != nil {
		return "", fmt.Errorf("error to gen jwt: %w", err)
	}
	return tokenString, nil
}

func (u *User) scan(rows *sql.Rows) (entity.User, error) {
	var user entity.User
	err := rows.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("error to scan: %w", err)
	}
	return user, nil
}

func (u *User) Register(ctx context.Context, db *sql.DB, form *dtos.RegisterForm) (string, error) {
	user, err := uow.NewTransaction(db, func() (*entity.User, error) {
		user, err := entity.NewUser(form.Username, form.Password)
		if err != nil {
			return nil, err
		}
		_, err = db.ExecContext(ctx, "INSERT INTO users(id, username, password) VALUES (?,?,?)",
			user.Id, user.Username, user.Password)
		return user, err
	}).Exec()

	if err != nil {
		return "", err
	}

	return u.genJwt(user)
}

func (u *User) Login(ctx context.Context, db *sql.DB, form *dtos.LoginForm) (string, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, username, password FROM users WHERE username = ?", form.Username)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error to get user by username: %w", err)
	}
	defer rows.Close()

	rows.Next()
	user, err := u.scan(rows)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error to get user by username: %w", err)
	}
	if err := user.ComparePassword(form.Password); err != nil {
		fmt.Println(err)
		return "", errors.New("invalid username or password")
	}
	return u.genJwt(&user)
}
