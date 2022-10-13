package entity

import (
	"HolyCrusade/internal/entity/models"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type Consumer struct {
	App      core.Application
	Handlers map[string]func([]byte)
}

type Handler struct {
	App               core.Application
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
			go func(hand func([]byte), value []byte) {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Failed to handle message from MQ")
						return
					}
				}()
				hand(value)
			}(h, msg.Value)

		}
	}
}

func (h *Handler) NewUser(value []byte) {
	var nu core.NewUser
	err := json.Unmarshal(value, &nu)
	if err != nil {
		log.Println("Can't unmarshal the byte array")
		return
	}

	tx, err := h.App.DB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to start transaction")
		return
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(context.Background())
			if err != nil {
				log.Println("Failed to rollback transaction")
				return
			}
		} else {
			err := tx.Commit(context.Background())
			if err != nil {
				log.Println("Failed to commit transaction")
				return
			}
		}
	}()

	var u models.User
	u.Token = nu.UserToken
	uID, err := h.UserRepository.Insert(u)
	if err != nil {
		log.Println(err)
		return
	}

	c := models.City{
		UserID: uID,
		Name:   fmt.Sprintf("City of User %d", uID),
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
	}
}
