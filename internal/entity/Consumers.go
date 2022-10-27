package entity

import (
	"HolyCrusade/internal/entity/models"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mitchellh/mapstructure"
	"log"
)

type Consumer struct {
	Handlers map[string]func(interface{}) error
}

type Handler struct {
	UserRepository    repository.UserRepository
	CityRepository    repository.CityRepository
	BalanceRepository repository.BalanceRepository
}

func (c *Consumer) ListenMQ() error {
	app := core.GetApp()

	for {
		msg, err := app.MQ.Reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var mqe core.MQEvent
		err = json.Unmarshal(msg.Value, &mqe)
		if err != nil {
			log.Println("Can't unmarshal the byte array")
			return err
		}

		if h, ok := c.Handlers[mqe.Type]; ok {
			go func(hand func(interface{}) error, value interface{}) {
				err := hand(value)
				if err != nil {
					log.Println("Failed to consume job")
				}
			}(h, mqe.Payload)
		}
	}
}

func (h *Handler) NewUser(value interface{}) error {
	app := core.GetApp()
	var nu core.NewUser
	err := mapstructure.Decode(value, &nu)
	if err != nil {
		log.Println("Can't decode value into right struct")
		return err
	}

	tx, err := app.DB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to start transaction")
		return err
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
	u.ChatID = nu.ChatID
	uID, err := h.UserRepository.Insert(context.Background(), u)
	if err != nil {
		log.Println(err)
		return err
	}

	c := models.City{
		UserID: uID,
		Name:   fmt.Sprintf("City of User %d", uID),
		Rating: 0,
	}
	cID, err := h.CityRepository.Insert(context.Background(), c)
	if err != nil {
		log.Println(err)
		return err
	}

	b := models.Balance{CityID: cID}

	_, err = h.BalanceRepository.Insert(context.Background(), b)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
