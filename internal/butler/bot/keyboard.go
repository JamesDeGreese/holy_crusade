package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("😎 +1 warrior", "add_warrior"),
		tgbotapi.NewInlineKeyboardButtonData("🥸 +1 peasant", "add_peasant"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💯 get statistics", "show_stats"),
	),
)
