package main

import (
	"HolyCrusade/internal/entity"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
)

func main() {
	a := core.InitApp("config/entity_service.yml").WithDB().WithKafka()

	ur := repository.UserRepository{}.Init(a.DB)
	cr := repository.CityRepository{}.Init(a.DB)
	br := repository.BalanceRepository{}.Init(a.DB)

	var h = entity.Handler{
		UserRepository:    ur,
		CityRepository:    cr,
		BalanceRepository: br,
	}

	c := entity.Consumer{
		Handlers: map[string]func([]byte) error{
			"new_user": h.NewUser,
		},
	}

	err := c.ListenMQ()
	if err != nil {
		panic(err)
	}
}
