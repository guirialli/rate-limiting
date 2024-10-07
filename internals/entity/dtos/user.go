package dtos

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
