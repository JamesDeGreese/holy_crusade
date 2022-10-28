package main

import (
	"HolyCrusade/internal/butler/bot"
	"HolyCrusade/pkg/core"
)

func main() {
	a := core.InitApp("config/butler.yml").WithTelegramBot()

	b := bot.GameBot{}
	b.SetBot(a.TgBot)

	b.Start()
}
