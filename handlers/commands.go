package handlers

import (
	"fmt"
	"log"
	"strings"

	"vps-monitor-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	Telegram *services.TelegramService
}

func NewCommandHandler(telegram *services.TelegramService) *CommandHandler {
	return &CommandHandler{
		Telegram: telegram,
	}
}

func (h *CommandHandler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s (ChatID: %d)", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

	h.Telegram.AddChatID(update.Message.Chat.ID)

	switch {
	case strings.HasPrefix(update.Message.Text, "/start"):
		h.handleStart(update.Message)
	case strings.HasPrefix(update.Message.Text, "/usage"):
		h.handleUsage(update.Message)
	case strings.HasPrefix(update.Message.Text, "/help"):
		h.handleHelp(update.Message)
	case strings.HasPrefix(update.Message.Text, "/status"):
		h.handleStatus(update.Message)
	case strings.HasPrefix(update.Message.Text, "/top"):
		h.handleTopProcesses(update.Message)
	}
}

func (h *CommandHandler) handleStart(msg *tgbotapi.Message) {
	response := fmt.Sprintf(
		"👋 *Selamat datang di VPS Monitor Bot!*\n\n" +
			"Bot ini akan memonitor VPS Anda dan mengirimkan alert secara otomatis ke chat ini jika:\n" +
			"• CPU usage ≥ 90%%\n\n" +
			"✅ *Status:* Anda sudah terdaftar untuk menerima alert di chat ini.\n\n" +
			"*Perintah yang tersedia:*\n" +
			"/usage - Cek penggunaan CPU, RAM, dan Disk\n" +
			"/status - Status lengkap sistem\n" +
			"/top - 5 proses dengan CPU tertinggi\n" +
			"/help - Bantuan",
	)
	h.Telegram.SendMessage(msg.Chat.ID, response)
}

func (h *CommandHandler) handleUsage(msg *tgbotapi.Message) {
	stats, err := services.GetSystemStats()
	if err != nil {
		h.Telegram.SendMessage(msg.Chat.ID, "❌ Error membaca statistik sistem")
		return
	}

	h.Telegram.SendMessage(msg.Chat.ID, services.FormatStats(stats))
}

func (h *CommandHandler) handleStatus(msg *tgbotapi.Message) {
	stats, err := services.GetSystemStats()
	if err != nil {
		h.Telegram.SendMessage(msg.Chat.ID, "❌ Error membaca status sistem")
		return
	}

	status := "✅ *Sistem Normal*"
	if stats.CPUUsage >= 90 {
		status = "🔥 *CPU Critical*"
	} else if stats.CPUUsage >= 70 {
		status = "⚠️ *CPU Warning*"
	}

	response := fmt.Sprintf("%s\n\n%s", status, services.FormatStats(stats))
	h.Telegram.SendMessage(msg.Chat.ID, response)
}

func (h *CommandHandler) handleTopProcesses(msg *tgbotapi.Message) {
	response := "🔧 *Top 5 Processes by CPU*\n\n" +
		"Fitur ini memerlukan implementasi dengan `ps` command.\n" +
		"Gunakan: `ps aux --sort=-%cpu | head -n 6`"
	h.Telegram.SendMessage(msg.Chat.ID, response)
}

func (h *CommandHandler) handleHelp(msg *tgbotapi.Message) {
	response := "📖 *Bantuan VPS Monitor Bot*\n\n" +
		"*Perintah:*\n" +
		"/start - Memulai bot dan mengaktifkan alert\n" +
		"/usage - Statistik penggunaan resource (CPU, RAM, Disk)\n" +
		"/status - Status kesehatan sistem lengkap\n" +
		"/top - Daftar proses dengan CPU tertinggi\n" +
		"/help - Menampilkan pesan ini\n\n" +
		"*Fitur Monitoring:*\n" +
		"• Monitoring CPU otomatis setiap 10 detik\n" +
		"• Alert otomatis ketika CPU ≥ 90%\n" +
		"• Monitoring Memory dan Disk\n" +
		"• Informasi Network I/O\n\n" +
		"*Catatan:*\n" +
		"Setiap user yang mengirim pesan ke bot ini akan otomatis menerima alert monitoring."
	h.Telegram.SendMessage(msg.Chat.ID, response)
}