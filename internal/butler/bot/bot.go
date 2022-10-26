package bot

import (
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

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := gb.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
