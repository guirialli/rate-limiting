//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
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

func NewAuthController(db *sql.DB, user usecases.IUser) *controller.Auth {
	wire.Build(
		controller.NewAuth,
	)
	return &controller.Auth{}
}
