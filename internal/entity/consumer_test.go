package entity

import (
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"testing"
)

func initHandler() Handler {
	a := core.InitApp("../../config/entity_service.yml").WithDB().WithKafka()

	ur := repository.UserRepository{}.Init(a.DB)
	cr := repository.CityRepository{}.Init(a.DB)
	br := repository.BalanceRepository{}.Init(a.DB)

	var h = Handler{
		UserRepository:    ur,
		CityRepository:    cr,
		BalanceRepository: br,
	}

	return h
}

func TestNewUserHandler(t *testing.T) {
	h := initHandler()

	nu := core.NewUser{ChatID: 1234567890}

	value, err := json.Marshal(nu)
	if err != nil {
		t.Fail()
	}

	err = h.NewUser(value)
	if err != nil {
		t.Fail()
	}

	u, err := h.UserRepository.GetByChatId(context.Background(), 1234567890)

	if u.ID == 0 || err != nil {
		t.Fail()
	}

	c, err := h.CityRepository.GetByUserID(context.Background(), u.ID)

	if c.ID == 0 || err != nil {
		t.Fail()
	}

	b, err := h.BalanceRepository.GetByCityID(context.Background(), c.ID)

	if b.ID == 0 || err != nil {
		t.Fail()
	}

	err = h.UserRepository.Delete(context.Background(), u.ID)

	if err != nil {
		t.Fail()
	}
}
