package core

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Application struct {
	Config Config
	MQ     *kafka.Conn
	DB     *pgx.Conn
}

type Config struct {
	Version  string   `yaml:"version"`
	Kafka    KafkaCfg `yaml:"kafka"`
	Database Database `yaml:"database"`
}

type KafkaCfg struct {
	Topic     string `yaml:"topic"`
	Partition int    `yaml:"partition"`
	Address   string `yaml:"address"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	UserName string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"dbName"`
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

func (a *Application) WithKafka() (*Application, error) {
	topic := a.Config.Kafka.Topic
	partition := a.Config.Kafka.Partition
	address := a.Config.Kafka.Address
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	a.MQ = conn

	return a, err
}

func (a *Application) WithDB() (*Application, error) {
	dbURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		a.Config.Database.Driver,
		a.Config.Database.UserName,
		a.Config.Database.Password,
		a.Config.Database.Host,
		a.Config.Database.Port,
		a.Config.Database.Database,
	)

	DBConn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect with database:", err)
	}

	a.DB = DBConn

	defer func(DBConn *pgx.Conn, ctx context.Context) {
		err := DBConn.Close(ctx)
		if err != nil {
			log.Fatal("failed to close database connection:", err)
		}
	}(DBConn, context.Background())

	return a, err
}
