package handlers

import (
	"log"
	"time"
	"vps-monitor-bot/config"
	"vps-monitor-bot/services"
)

type Monitor struct {
	Config        *config.Config
	Telegram      *services.TelegramService
	lastAlert     time.Time
	alertCooldown time.Duration
}

func NewMonitor(cfg *config.Config, telegram *services.TelegramService) *Monitor {
	return &Monitor{
		Config:        cfg,
		Telegram:      telegram,
		alertCooldown: 5 * time.Minute,
	}
}

func (m *Monitor) Start() {
	log.Println("🚀 Memulai CPU monitoring...")

	ticker := time.NewTicker(time.Duration(m.Config.MonitorInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.checkCPU()
	}
}

func (m *Monitor) checkCPU() {
	cpuUsage, err := services.GetCPUUsage()
	if err != nil {
		log.Printf("Error membaca CPU: %v", err)
		return
	}

	log.Printf("CPU Usage: %.2f%%", cpuUsage)

	if cpuUsage >= m.Config.CPUThreshold {
		if time.Since(m.lastAlert) > m.alertCooldown {
			stats, err := services.GetSystemStats()
			if err != nil {
				log.Printf("Error membaca system stats: %v", err)
				return
			}

			m.Telegram.SendAlert(cpuUsage, stats)
			m.lastAlert = time.Now()
			log.Printf("🚨 Alert terkirim! CPU: %.2f%%", cpuUsage)
		} else {
			log.Printf("⏳ Alert di-skip (cooldown aktif)")
		}
	}
}