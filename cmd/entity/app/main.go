package main

import (
	"HolyCrusade/internal/entity"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"fmt"
	"log"
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

	done := make(chan bool)

	go func(ch chan bool) {
		err := entity.ListenMQ("new_user", h.NewUser)
		if err != nil {
			ch <- true
			log.Panic(err)
		}
	}(done)

	go func(ch chan bool) {
		err := entity.ListenMQ("city_info_req", h.CityInfo)
		if err != nil {
			ch <- true
			log.Panic(err)
		}
	}(done)

	go func(ch chan bool) {
		err := entity.ListenMQ("add_worker_req", h.AddWorker)
		if err != nil {
			ch <- true
			log.Panic(err)
		}
	}(done)

	go func(ch chan bool) {
		err := entity.ListenMQ("add_solder_req", h.AddSolder)
		if err != nil {
			ch <- true
			log.Panic(err)
		}
	}(done)

	select {
	case <-done:
		fmt.Println("One of readers down, done")
		return
	}
}
