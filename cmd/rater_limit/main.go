package main

import "github.com/guirialli/rater_limit/internals/infra/database"

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
}
