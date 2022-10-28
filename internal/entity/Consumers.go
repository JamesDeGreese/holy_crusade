package entity

import (
	"HolyCrusade/internal/entity/models"
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"log"
)

type Handler struct {
	UserRepository    repository.UserRepository
	CityRepository    repository.CityRepository
	BalanceRepository repository.BalanceRepository
}

func ListenMQ(topic string, handler func(interface{}) error) error {
	app := core.GetApp()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{app.Config.Kafka.Address},
		Topic:   topic,
		GroupID: "default_group",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var mqe interface{}
		err = json.Unmarshal(msg.Value, &mqe)
		if err != nil {
			log.Println("Can't unmarshal the bytes array")
			return err
		}

		go func(hand func(interface{}) error, value interface{}) {
			err := hand(value)
			if err != nil {
				log.Println("Failed to consume job")
			}
		}(handler, mqe)
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

	if exu, _ := h.UserRepository.GetByChatId(context.Background(), nu.ChatID); exu.ID != 0 {
		log.Println("user already exists")
		return errors.New("user already exists")
	}

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

func (h *Handler) CityInfo(value interface{}) error {
	var ci core.CityInfoReq
	err := mapstructure.Decode(value, &ci)
	if err != nil {
		log.Println("Can't decode value into right struct")
		return err
	}

	user, err := h.UserRepository.GetByChatId(context.Background(), ci.ChatID)
	if err != nil || user.ID == 0 {
		log.Println("Can't get user")
		return err
	}

	city, err := h.CityRepository.GetByUserID(context.Background(), user.ID)
	if err != nil || user.ID == 0 {
		log.Println("Can't get city")
		return err
	}

	balance, err := h.BalanceRepository.GetByCityID(context.Background(), city.ID)
	if err != nil || user.ID == 0 {
		log.Println("Can't get balance")
		return err
	}

	err = CityInfoResponse(user.ChatID, city, balance)
	if err != nil {
		return err
	}

	return nil
}
