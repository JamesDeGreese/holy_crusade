package main

import (
	"HolyCrusade/pkg/core"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	app := core.InitApp("config/entity_service.yml").WithDB()

	m, err := migrate.New(
		"file://internal/entity/migrations",
		app.DB.Config().ConnString())
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
