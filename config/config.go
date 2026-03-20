package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken        string
	AllowedChatIDs  []int64
	CPUThreshold    float64
	MonitorInterval int
}

func Load() *Config {
	threshold, _ := strconv.ParseFloat(getEnv("CPU_THRESHOLD", "90"), 64)
	interval, _ := strconv.Atoi(getEnv("MONITOR_INTERVAL", "10"))

	var chatIDs []int64
	if idsStr := os.Getenv("ALLOWED_CHAT_IDS"); idsStr != "" {
		parts := strings.Split(idsStr, ",")
		for _, part := range parts {
			if id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64); err == nil {
				chatIDs = append(chatIDs, id)
			}
		}
	}

	return &Config{
		BotToken:        getEnv("BOT_TOKEN", ""),
		AllowedChatIDs:  chatIDs,
		CPUThreshold:    threshold,
		MonitorInterval: interval,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}