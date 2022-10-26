package core

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var app *Application = nil

type Application struct {
	Config Config
	MQ     MQ
	DB     *pgx.Conn
	TgBot  *tgbotapi.BotAPI
}

type MQ struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

type Config struct {
	Version  string   `yaml:"version"`
	Kafka    KafkaCfg `yaml:"kafka"`
	Database Database `yaml:"database"`
	TgBot    TgCfg    `yaml:"telegram_bot"`
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

type TgCfg struct {
	Token string `yaml:"token"`
	Debug bool   `yaml:"debug"`
}

func GetApp() *Application {
	if app == nil {
		app = &Application{}
	}
	return app
}

func (c *Config) ParseFile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(contents, c)
}

func InitApp(cfgPath string) *Application {
	var cfg Config

	err := cfg.ParseFile(cfgPath)
	if err != nil {
		panic(err)
	}

	a := GetApp()
	a.Config = cfg

	return a
}

func (a *Application) WithKafka() *Application {
	var mq MQ

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{a.Config.Kafka.Address},
		Topic:   a.Config.Kafka.Topic,
		GroupID: "default_group",
	})

	w := &kafka.Writer{
		Addr:     kafka.TCP(a.Config.Kafka.Address),
		Topic:    a.Config.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	mq.Reader = r
	mq.Writer = w

	a.MQ = mq

	return a
}

func (a *Application) WithDB() *Application {
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

	return a
}

func (a *Application) WithTelegramBot() *Application {
	bot, err := tgbotapi.NewBotAPI(a.Config.TgBot.Token)
	if err != nil {
		log.Fatal("failed to run tgbot:", err)
	}

	bot.Debug = a.Config.TgBot.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	a.TgBot = bot

	return a
}
