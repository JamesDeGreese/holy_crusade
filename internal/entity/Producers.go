package entity

import (
	"HolyCrusade/internal/entity/models"
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

func CityInfoResponse(chatID int64, c models.City, b models.Balance) error {
	app := core.GetApp()

	writer := kafka.Writer{
		Addr:     kafka.TCP(app.Config.Kafka.Address),
		Topic:    "response",
		Balancer: &kafka.LeastBytes{},
	}
	uuidv4, _ := uuid.NewV4()
	v, _ := json.Marshal(
		core.Response{
			ChatID: chatID,
			Payload: core.CityInfoRes{
				Name:       c.Name,
				Rating:     c.Rating,
				Gold:       b.Gold,
				Population: b.Population,
				Workers:    b.Workers,
				Solders:    b.Solders,
				Heroes:     b.Heroes,
			},
		})
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   uuidv4.Bytes(),
			Value: v,
		},
	)

	return err
}
