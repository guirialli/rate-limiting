package vos

type LoginForm struct {
	username string `form:"username"`
	password string `form:"password"`
}

type RegisterForm struct {
	username string `form:"username"`
	password string `form:"password"`
}
