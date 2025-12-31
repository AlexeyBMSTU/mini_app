package bot

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/user"
	
	_ "github.com/lib/pq"
)

type Bot struct {
	API     *tgbotapi.BotAPI
	Config  *config.Config
	DB      *sql.DB
	Updates tgbotapi.UpdatesChannel
	UserRepo *user.SQLRepository
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
	if err := b.initDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	
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

func (b *Bot) initDB() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		b.Config.PostgresUser, b.Config.PostgresPassword, b.Config.PostgresHost, b.Config.PostgresPort, b.Config.PostgresDB)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	
	b.DB = db
	
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	
	log.Println("✅ Bot connected to database")
	
	b.UserRepo = user.NewSQLRepository(db)
	
	err = b.UserRepo.CreateTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}
	
	log.Println("✅ Bot database tables created")
	
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