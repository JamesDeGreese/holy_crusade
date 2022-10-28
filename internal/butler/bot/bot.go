package bot

import (
	"HolyCrusade/pkg/core"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GameBot struct {
	bot *tgbotapi.BotAPI
}

func (gb *GameBot) SetBot(bot *tgbotapi.BotAPI) {
	gb.bot = bot
}

func (gb *GameBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := gb.bot.GetUpdatesChan(u)
	uuidv4, _ := uuid.NewV4()

	app := core.GetApp()

	go func() {
		err := func() error {
			app := core.GetApp()
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{app.Config.Kafka.Address},
				Topic:   "response",
				GroupID: "default_group",
			})
			for {
				msg, err := reader.ReadMessage(context.Background())
				if err != nil {
					return err
				}

				var res core.Response
				err = json.Unmarshal(msg.Value, &res)
				if err != nil {
					return err
				}

				text, err := json.MarshalIndent(res.Payload, "", "  ")
				if err != nil {
					return err
				}
				m := tgbotapi.NewMessage(res.ChatID, string(text))
				if _, err := gb.bot.Send(m); err != nil {
					log.Panic(err)
				}
			}

			return nil
		}()
		if err != nil {
			log.Println("Failed to process response")
		}
	}()

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		switch update.Message.Text {
		case "/start":
			writer := kafka.Writer{
				Addr:     kafka.TCP(app.Config.Kafka.Address),
				Topic:    "new_user",
				Balancer: &kafka.LeastBytes{},
			}
			v, _ := json.Marshal(core.NewUser{ChatID: update.Message.Chat.ID})
			err := writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   uuidv4.Bytes(),
					Value: v,
				},
			)
			if err != nil {
				log.Panic(err)
			}
		case "/city":
			writer := kafka.Writer{
				Addr:     kafka.TCP(app.Config.Kafka.Address),
				Topic:    "city_info_req",
				Balancer: &kafka.LeastBytes{},
			}
			v, _ := json.Marshal(core.CityInfoReq{ChatID: update.Message.Chat.ID})
			err := writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   uuidv4.Bytes(),
					Value: v,
				},
			)
			if err != nil {
				log.Panic(err)
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
			if _, err := gb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
