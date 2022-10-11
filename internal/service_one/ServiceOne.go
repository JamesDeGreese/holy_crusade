package service_one

import (
	"HolyCrusade/pkg/core"
	"github.com/segmentio/kafka-go"
	"log"
)

type ServiceOne struct {
	App core.Application
}

func (s ServiceOne) SendToQueue() {
	_, err := s.App.MQ.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	println("Done")
}
