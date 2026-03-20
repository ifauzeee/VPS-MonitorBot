package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"vps-monitor-bot/config"
	"vps-monitor-bot/handlers"
	"vps-monitor-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	// Load konfigurasi
	cfg := config.Load()

	if cfg.BotToken == "" {
		log.Fatal("❌ BOT_TOKEN harus di-set di environment variable")
	}

	// Inisialisasi Telegram service
	telegramService, err := services.NewTelegramService(cfg.BotToken, cfg.AllowedChatIDs)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi Telegram bot: %v", err)
	}

	// Setup command handler
	commandHandler := handlers.NewCommandHandler(telegramService)

	// Setup dan jalankan monitor di goroutine terpisah
	monitor := handlers.NewMonitor(cfg, telegramService)
	go monitor.Start()

	// Setup update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := telegramService.Bot.GetUpdatesChan(u)

	// Channel untuk graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop
	go func() {
		for update := range updates {
			commandHandler.HandleUpdate(update)
		}
	}()

	log.Println("✅ Bot berjalan. Tekan Ctrl+C untuk menghentikan.")
	<-sigChan

	log.Println("👋 Mematikan bot...")
}
