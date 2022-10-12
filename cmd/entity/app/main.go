package main

import (
	"HolyCrusade/internal/entity"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
)

func main() {
	var a core.Application
	a.Init("config/entity_service.yml").WithDB().WithKafka()

	h := entity.Handler{
		UserRepository:    repository.UserRepository{DB: a.DB},
		CityRepository:    repository.CityRepository{DB: a.DB},
		BalanceRepository: repository.BalanceRepository{DB: a.DB},
	}

	c := entity.Consumer{
		App: a,
		Handlers: map[string]func([]byte){
			"new_user": h.NewUser,
		},
	}

	c.ListenMQ()
}
