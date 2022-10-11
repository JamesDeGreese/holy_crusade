package core

import (
	"context"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Application struct {
	Config Config
	MQ     *kafka.Conn
}

type Config struct {
	Version string   `yaml:"version"`
	Kafka   KafkaCfg `yaml:"kafka"`
}

type KafkaCfg struct {
	Topic     string `yaml:"topic"`
	Partition int    `yaml:"partition"`
	Address   string `yaml:"address"`
}

func (c *Config) ParseFile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(contents, c)
}

func (a *Application) Init(cfgPath string) *Application {
	var cfg Config

	err := cfg.ParseFile(cfgPath)
	if err != nil {
		panic(err)
	}

	a.Config = cfg

	return a
}

func (a *Application) WithKafka() *Application {
	topic := a.Config.Kafka.Topic
	partition := a.Config.Kafka.Partition
	address := a.Config.Kafka.Address
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	a.MQ = conn

	return a
}
