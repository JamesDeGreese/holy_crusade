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

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		switch update.Message.Text {
		case "/start":
			v, _ := json.Marshal(
				core.MQEvent{
					Type:    "new_user",
					Payload: core.NewUser{ChatID: update.Message.Chat.ID},
				})
			err := app.MQ.Writer.WriteMessages(context.Background(),
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
