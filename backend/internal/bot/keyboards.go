package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateStartKeyboard(botUsername string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"🚀 Открыть Servatory",
				fmt.Sprintf("https://t.me/%s?startapp=webapp", botUsername),
			),
		),
	)
}
