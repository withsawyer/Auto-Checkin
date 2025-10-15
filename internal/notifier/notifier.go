package notifier

import (
	"auto-checkin/internal/config"
	"log"
)

type Notifier struct {
	wecomWebhook     string
	telegramBotToken string
	telegramChatID   string
}

func New(cfg config.Notifications) *Notifier {
	return &Notifier{
		wecomWebhook:     cfg.WeCom.Webhook,
		telegramBotToken: cfg.Telegram.BotToken,
		telegramChatID:   cfg.Telegram.ChatID,
	}
}

func (n *Notifier) SendWeCom(message string) error {
	log.Printf("企微消息推送: %s", message)
	return nil
}

func (n *Notifier) SendTelegram(message string) error {
	log.Printf("Telegram消息推送: %s", message)
	return nil
}
