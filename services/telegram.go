package services

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	Bot     *tgbotapi.BotAPI
	ChatIDs map[int64]bool
	mu      sync.RWMutex
}

func NewTelegramService(token string, initialChatIDs []int64) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	chatIDsMap := make(map[int64]bool)
	for _, id := range initialChatIDs {
		chatIDsMap[id] = true
	}

	return &TelegramService{
		Bot:     bot,
		ChatIDs: chatIDsMap,
	}, nil
}

func (t *TelegramService) AddChatID(chatID int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.ChatIDs[chatID] {
		t.ChatIDs[chatID] = true
		log.Printf("New chat session added: %d", chatID)
	}
}

func (t *TelegramService) SendAlert(cpuUsage float64, stats *SystemStats) {
	message := fmt.Sprintf(
		"🚨 *ALERT: High CPU Usage Detected!*\n\n"+
			"🔥 Current CPU Usage: `%.2f%%`\n"+
			"⚠️ Threshold: `90%%`\n\n"+
			"%s",
		cpuUsage,
		FormatStats(stats),
	)

	t.mu.RLock()
	ids := make([]int64, 0, len(t.ChatIDs))
	for id := range t.ChatIDs {
		ids = append(ids, id)
	}
	t.mu.RUnlock()

	for _, chatID := range ids {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "Markdown"
		if _, err := t.Bot.Send(msg); err != nil {
			log.Printf("Error sending alert to %d: %v", chatID, err)
		}
	}
}

func (t *TelegramService) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := t.Bot.Send(msg)
	return err
}