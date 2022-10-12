package entity

import (
	"HolyCrusade/internal/entity/models"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"log"
)

type Consumer struct {
	App      core.Application
	Handlers map[string]func([]byte)
}

type Handler struct {
	UserRepository    repository.UserRepository
	CityRepository    repository.CityRepository
	BalanceRepository repository.BalanceRepository
}

func (c *Consumer) ListenMQ() {
	for {
		msg, err := c.App.MQ.Reader.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}

		if h, ok := c.Handlers[string(msg.Key)]; ok {
			go h(msg.Value)
		}
	}
}

// TODO: Тесты для этого хендлера

func (h *Handler) NewUser(value []byte) {
	var nu core.NewUser
	err := json.Unmarshal(value, &nu)
	if err != nil {
		log.Println("Can't unmarshal the byte array")
		return
	}

	var u models.User
	u.Token = nu.UserToken
	uID, err := h.UserRepository.Insert(u)
	if err != nil {
		log.Println(err)
		return
	}

	c := models.City{
		UserID: uID,
		Name:   "",
		Rating: 0,
	}
	cID, err := h.CityRepository.Insert(c)
	if err != nil {
		log.Println(err)
		return
	}

	b := models.Balance{
		CityID:     cID,
		Gold:       0,
		Population: 0,
		Workers:    0,
		Solders:    0,
		Heroes:     0,
	}

	_, err = h.BalanceRepository.Insert(b)

	if err != nil {
		log.Println(err)
		return
	}

}
