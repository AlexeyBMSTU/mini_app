package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handleStartCommand(message)
	case "help":
		b.handleHelpCommand(message)
	case "about":
		b.handleAboutCommand(message)
	case "time":
		b.handleTimeCommand(message)
	default:
		b.handleUnknownCommand(message)
	}
}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	log.Printf("📩 Received message from @%s: %s", message.From.UserName, message.Text)

	text := strings.ToLower(strings.TrimSpace(message.Text))

	switch {
	case strings.Contains(text, "привет") || strings.Contains(text, "здравствуй"):
		b.handleGreeting(message)
	case strings.Contains(text, "как дела") || strings.Contains(text, "как ты"):
		b.handleHowAreYou(message)
	case strings.Contains(text, "спасибо"):
		b.handleThanks(message)
	case strings.Contains(text, "пока") || strings.Contains(text, "до свидания"):
		b.handleGoodbye(message)
	case strings.Contains(text, "что ты умеешь") || strings.Contains(text, "функционал"):
		b.handleCapabilities(message)
	case strings.Contains(text, "новости") || strings.Contains(text, "обновления"):
		b.handleNews(message)
	case strings.Contains(text, "помощь"):
		b.handleHelpCommand(message)
	default:
		b.handleDefaultResponse(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := fmt.Sprintf(`
👋 Привет, %s!
✨ <b>Это первая версия бота %s</b>

Я помогу тебе с... [добавьте описание функционала]
Просто напиши мне что-нибудь или используй команды!`,
		message.From.FirstName,
		b.API.Self.UserName)

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = CreateStartKeyboard(b.API.Self.UserName)

	if _, err := b.API.Send(msg); err != nil {
		log.Printf("Failed send message: %v", err)
	} else {
		log.Printf("✅ Send message @%s", message.From.UserName)
	}
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `📚 <b>Доступные команды:</b>

/start - Начать работу с ботом
/help - Показать это сообщение
/about - О боте
/time - Текущее время

🤖 <b>Также я понимаю:</b>
• Приветствия (привет, здравствуйте)
• Вопросы о делах (как дела?)
• Благодарности (спасибо)
• Прощания (пока, до свидания)
• Вопросы о возможностях (что ты умеешь?)
• Запрос новостей (новости, обновления)

🔗 <b>Мини-приложение:</b>
Используйте кнопку "Открыть Servatory" для запуска мини-приложения.`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Все команды", "all_commands"),
		),
	)
	b.API.Send(msg)
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID,
		"❌ Неизвестная команда. Используйте /help для списка команд.\n\n"+
			"Или просто напишите мне сообщение — я постараюсь понять!")
	b.API.Send(msg)
}

func (b *Bot) handleAboutCommand(message *tgbotapi.Message) {
	aboutText := fmt.Sprintf(`🤖 <b>О боте %s</b>

Версия: 1.0.0
Разработчик: [Ваше имя/компания]
Дата создания: %s

Этот бот создан для... [описание цели бота]
Исходный код: [ссылка на GitHub, если есть]`,
		b.API.Self.UserName,
		time.Now().Format("02.01.2006"))

	msg := tgbotapi.NewMessage(message.Chat.ID, aboutText)
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}

func (b *Bot) handleTimeCommand(message *tgbotapi.Message) {
	currentTime := time.Now().Format("15:04:05 02.01.2006")
	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("⏰ Текущее время: <b>%s</b>", currentTime))
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}

func (b *Bot) handleGreeting(message *tgbotapi.Message) {
	responses := []string{
		fmt.Sprintf("👋 Привет, %s! Рад тебя видеть!", message.From.FirstName),
		fmt.Sprintf("Здравствуй, %s! Как твои дела?", message.From.FirstName),
		fmt.Sprintf("Приветствую, %s! Чем могу помочь?", message.From.FirstName),
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responses[time.Now().Unix()%int64(len(responses))])
	b.API.Send(msg)
}

func (b *Bot) handleHowAreYou(message *tgbotapi.Message) {
	responses := []string{
		"🤖 У меня всё отлично! Я просто программа, но стараюсь быть полезным!",
		"👍 Всё хорошо, спасибо! Готов помочь тебе.",
		"✨ Отлично! Работаю в полную силу. А как ты?",
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responses[time.Now().Unix()%int64(len(responses))])
	b.API.Send(msg)
}

func (b *Bot) handleThanks(message *tgbotapi.Message) {
	responses := []string{
		"😊 Всегда рад помочь!",
		"🙏 Пожалуйста! Обращайся ещё.",
		"✨ Не за что! Буду рад помочь снова.",
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responses[time.Now().Unix()%int64(len(responses))])
	b.API.Send(msg)
}

func (b *Bot) handleGoodbye(message *tgbotapi.Message) {
	responses := []string{
		fmt.Sprintf("👋 До свидания, %s! Буду ждать нашего следующего общения!", message.From.FirstName),
		"Пока! Возвращайся скорее!",
		"До встречи! Не стесняйся обращаться снова!",
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responses[time.Now().Unix()%int64(len(responses))])
	b.API.Send(msg)
}

func (b *Bot) handleCapabilities(message *tgbotapi.Message) {
	capabilitiesText := `🤖 <b>Мои возможности:</b>

• Отвечать на команды (/start, /help, /about, /time)
• Поддерживать простой диалог
• Понимать базовые фразы (привет, пока, спасибо и т.д.)
• Предоставлять информацию о себе
• [Добавьте сюда ваш функционал]

📝 Для полного списка команд используй /help`

	msg := tgbotapi.NewMessage(message.Chat.ID, capabilitiesText)
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}

func (b *Bot) handleNews(message *tgbotapi.Message) {
	newsText := `📢 <b>Последние обновления:</b>

• <b>Версия 1.0.0</b> — Первый релиз бота
• Добавлены базовые команды
• Реализована обработка текстовых сообщений
• Создана структура для расширения функционала

🔮 <b>В планах:</b>
• [Добавьте планы по развитию]
• [Другой функционал]`

	msg := tgbotapi.NewMessage(message.Chat.ID, newsText)
	msg.ParseMode = "HTML"
	b.API.Send(msg)
}

func (b *Bot) handleDefaultResponse(message *tgbotapi.Message) {
	responses := []string{
		fmt.Sprintf("🤔 Извини, %s, я не совсем понял. Попробуй использовать команду /help или спроси что-то проще.", message.From.FirstName),
		"Прости, я ещё только учусь! Напиши /help чтобы узнать, что я умею.",
		fmt.Sprintf("%s, я пока не могу ответить на это. Используй команды или спроси что-то другое!", message.From.FirstName),
		"Интересный вопрос! Но мои возможности пока ограничены. Попробуй /help для списка команд.",
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responses[time.Now().Unix()%int64(len(responses))])
	b.API.Send(msg)
}
