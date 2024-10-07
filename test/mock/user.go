package mock

import (
	"github.com/guirialli/rater_limit/internals/entity/dtos"
)

type User struct{}

func NewUserMock() *User {
	return &User{}
}

func (u *User) RegisterForm() *dtos.RegisterForm {
	return &dtos.RegisterForm{
		Username: "test",
		Password: "T2@stt1231131",
	}
}

func (u *User) LoginForm() *dtos.LoginForm {
	return &dtos.LoginForm{
		Username: "test",
		Password: "T2@stt1231131",
	}
}
