package main

import (
	"fmt"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"github.com/guirialli/rater_limit/internals/infra/webserver/server"
)

// @title Rater Limit
// @version 1.0
// @description rater limit example
// @termsOfService: https://swagger.io/terms/

// @contact.name Guilherme Rialli
// @contact.url https://www.linkedin.com/in/guilherme-rialli-oliveira-1b826a150/
// @contact.email gui.rialli@gmail.com

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	fmt.Println("Start Rater Limit Application")
	db, err := database.NewMySql()
	if err != nil {
		panic(err)
	}

	err = db.TryConnection()
	if err != nil {
		panic(err)
	}

	if err = db.InitDatabase("init.sql"); err != nil {
		panic(err)
	}

	con, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	defer con.Close()

	cfg, err := config.LoadWebServerConfig()
	if err != nil {
		panic(err)
	}

	if err = server.NewServer(cfg, []router.UseRouter{
		NewBookRouter(con),
		NewAuthorRouter(con),
		NewAuthRouter(con),
		NewSwaggerRouter("./docs"),
	}).Start(NewRaterLimitMiddleware()); err != nil {
		panic(err)
	}
}
