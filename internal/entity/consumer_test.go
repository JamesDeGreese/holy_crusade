package entity

import (
	"HolyCrusade/internal/entity/repository"
	"HolyCrusade/pkg/core"
	"encoding/json"
	"testing"
)

func initHandler() Handler {
	var a core.Application
	a.Init("../../config/entity_service.yml").WithDB().WithKafka()

	h := Handler{
		App:               a,
		UserRepository:    repository.UserRepository{DB: a.DB},
		CityRepository:    repository.CityRepository{DB: a.DB},
		BalanceRepository: repository.BalanceRepository{DB: a.DB},
	}

	return h
}

func TestNewUserHandler(t *testing.T) {
	h := initHandler()

	nu := core.NewUser{UserToken: "TOKEN"}

	value, err := json.Marshal(nu)
	if err != nil {
		t.Fail()
	}

	h.NewUser(value)

	u, err := h.UserRepository.GetByToken("TOKEN")

	if u.ID == 0 || err != nil {
		t.Fail()
	}

	c, err := h.CityRepository.GetByUserID(u.ID)

	if c.ID == 0 || err != nil {
		t.Fail()
	}

	b, err := h.BalanceRepository.GetByCityID(c.ID)

	if b.ID == 0 || err != nil {
		t.Fail()
	}
}
