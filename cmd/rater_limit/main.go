package main

import (
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"github.com/guirialli/rater_limit/internals/infra/webserver/server"
)

func main() {
	db, err := database.NewMySql()
	if err != nil {
		panic(err)
	}
	if err := db.TryConnection(); err != nil {
		panic(err)
	}
	if err := db.InitDatabase("init.sql"); err != nil {
		panic(err)
	}

	con, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	cfg, err := config.LoadWebServerConfig()
	if err != nil {
		panic(err)
	}

	if err = server.NewServer(cfg, []router.UseRouter{
		NewBookRouter(con),
		NewAuthorRouter(con),
		NewAuthRouter(con),
	}).Start(); err != nil {
		panic(err)
	}
}
