//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/infra/webserver/controller"
	"github.com/guirialli/rater_limit/internals/infra/webserver/middleware"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"github.com/guirialli/rater_limit/internals/usecases"
)

// Use Cases Dependency
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
var setRaterLimitUseCaseDependency = wire.NewSet(
	newRateLimitUseCase,
	setUserUseCaseDependency,
	wire.Bind(new(usecases.IRaterLimit), new(*usecases.RaterLimit)),
)

var setHttpHandlerErrorDependency = wire.NewSet(
	controller.NewUtils,
	wire.Bind(new(controller.IHttpHandlerError), new(*controller.Utils)),
)

// controller dependency
var setAuthorControllerDependency = wire.NewSet(
	controller.NewAuthor,
	setBookUseCaseDependency,
	setAuthorUseCaseDependency,
	setHttpHandlerErrorDependency,
	wire.Bind(new(controller.IAuthor), new(*controller.Author)),
)

var setBookControllerDependency = wire.NewSet(
	controller.NewBook,
	setBookUseCaseDependency,
	setAuthorUseCaseDependency,
	setHttpHandlerErrorDependency,
	wire.Bind(new(controller.IBooks), new(*controller.Book)),
)

var setAuthControllerDependency = wire.NewSet(
	controller.NewAuth,
	setUserUseCaseDependency,
	wire.Bind(new(controller.IAuth), new(*controller.Auth)),
)

// router dependency
var setAuthMiddlewareDependency = wire.NewSet(
	newUser,
	wire.Bind(new(router.IAuthToken), new(*usecases.User)),
)

// DI
func NewAuthorController(db *sql.DB) *controller.Author {
	wire.Build(
		setAuthorUseCaseDependency,
		setBookUseCaseDependency,
		setHttpHandlerErrorDependency,
		controller.NewAuthor,
	)
	return &controller.Author{}
}

func NewBookController(db *sql.DB) *controller.Book {
	wire.Build(
		setAuthorUseCaseDependency,
		setBookUseCaseDependency,
		setHttpHandlerErrorDependency,
		controller.NewBook,
	)
	return &controller.Book{}
}

func NewAuthController(db *sql.DB) *controller.Auth {
	wire.Build(
		setUserUseCaseDependency,
		setHttpHandlerErrorDependency,
		controller.NewAuth,
	)
	return &controller.Auth{}
}

func NewAuthorRouter(db *sql.DB) *router.Author {
	wire.Build(
		setAuthorControllerDependency,
		setAuthMiddlewareDependency,
		router.NewAuthor,
	)
	return &router.Author{}
}

func NewBookRouter(db *sql.DB) *router.Book {
	wire.Build(
		setBookControllerDependency,
		setAuthMiddlewareDependency,
		router.NewBook,
	)
	return &router.Book{}
}

func NewAuthRouter(db *sql.DB) *router.Auth {
	wire.Build(
		setAuthControllerDependency,
		setHttpHandlerErrorDependency,
		router.NewAuth,
	)
	return &router.Auth{}
}

func NewRaterLimitMiddleware() *middleware.RaterLimit {
	wire.Build(
		setRaterLimitUseCaseDependency,
		middleware.NewRaterLimit,
	)
	return &middleware.RaterLimit{}
}

// utils constructors
// This function create a user use case without errors
func newUser() *usecases.User {
	jwtConfig := config.LoadJwtConfig()
	user, err := usecases.NewUser(jwtConfig.Secret, jwtConfig.ExpireIn, jwtConfig.UnitTime)
	if err != nil {
		panic(err)
	}
	return user
}

func newRateLimitUseCase(user usecases.IUser) *usecases.RaterLimit {
	cfg, _ := config.LoadRaterLimitConfig()
	rCfg := config.LoadRedisConfig()

	rdb := database.NewRedisClient[entity.RaterLimit](*rCfg)
	raterLimit, err := usecases.NewRaterLimit(user, *cfg, rdb)
	if err != nil {
		panic(err)
	}
	return raterLimit
}
