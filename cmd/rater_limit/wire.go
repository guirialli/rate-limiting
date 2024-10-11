//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/webserver/controller"
	"github.com/guirialli/rater_limit/internals/usecases"
)

var setAuthorUseCaseDependency = wire.NewSet(
	usecases.NewAuthor,
	wire.Bind(new(usecases.IAuthor), new(*usecases.Author)),
)

var setBookUseCaseDependency = wire.NewSet(
	usecases.NewBook,
	wire.Bind(new(usecases.IBook), new(*usecases.Book)),
)

var setUserUseCaseDependency = wire.NewSet(
	newUser,
	wire.Bind(new(usecases.IUser), new(*usecases.User)),
)

// This function create a user use case without errors
func newUser() *usecases.User {
	jwtConfig := config.LoadJwtConfig()
	user, err := usecases.NewUser(jwtConfig.Secret, jwtConfig.ExpireIn, jwtConfig.UnitTime)
	if err != nil {
		panic(err)
	}
	return user
}

func NewAuthorController(db *sql.DB) *controller.Author {
	wire.Build(
		setAuthorUseCaseDependency,
		setBookUseCaseDependency,
		controller.NewAuthor,
	)
	return &controller.Author{}
}

func NewBookController(db *sql.DB) *controller.Book {
	wire.Build(
		setAuthorUseCaseDependency,
		setBookUseCaseDependency,
		controller.NewBook,
	)
	return &controller.Book{}
}

func NewAuthController(db *sql.DB) *controller.Auth {
	wire.Build(
		setUserUseCaseDependency,
		controller.NewAuth,
	)
	return &controller.Auth{}
}
