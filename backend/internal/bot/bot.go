package bot

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mini-app-backend/internal/config"
)

type Bot struct {
	API     *tgbotapi.BotAPI
	Config  *config.Config
	Updates tgbotapi.UpdatesChannel
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		API:    api,
		Config: cfg,
	}

	if err := bot.setup(); err != nil {
		return nil, err
	}

	return bot, nil
}

func (b *Bot) setup() error {
	_, err := b.API.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	if err != nil {
		log.Printf("Failed delete webhook: %v", err)
	}

	time.Sleep(2 * time.Second)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.Updates = b.API.GetUpdatesChan(u)

	log.Printf("🤖 Auth as %s", b.API.Self.UserName)
	return nil
}

func (b *Bot) Start() {
	log.Println("✅ Bot started!")

	for update := range b.Updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.HandleCommand(update.Message)
			continue
		}

		if update.Message.Text != "" {
			b.HandleMessage(update.Message)
			continue
		}

	}
}