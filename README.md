# VPS Monitor Telegram Bot

A Telegram bot to monitor VPS CPU load and system statistics in real-time, written in Go.

## Features
- 📊 **Real-time Monitoring**: Monitor CPU, RAM, Disk, Load Average, and Uptime.
- 🚨 **Dynamic CPU Alerts**: Automatically send notifications to Telegram if CPU usage exceeds a threshold (default 90%). 
- 🤖 **Auto-Subscription**: Every user who interacts with the bot will automatically receive alerts.
- 🌐 **Network I/O**: Track network throughput (upload/download).
- ⚙️ **Customizable**: Adjustable CPU threshold and monitoring interval.

## Installation

### 1. Install Go
If Go is not yet installed on your system:
```bash
sudo apt update
sudo apt install -y golang-go
```

### 2. Setup Configuration
Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
Edit `.env` and provide your `BOT_TOKEN`. `ALLOWED_CHAT_IDS` is now optional (auto-subscribed on chat).

### 3. Build & Run
```bash
go mod tidy
go build -o vps-monitor-bot main.go
./vps-monitor-bot
```

## How to Use
1. Start the bot on Telegram.
2. Send `/start` or any message.
3. The bot will automatically add you to the alert list.

## Bot Commands
| Command | Description |
|---------|-----------|
| `/start` | Initialize the bot and activate alerts for your chat session |
| `/usage` | Resource usage statistics (CPU, RAM, Disk) |
| `/status` | Full system health status |
| `/top` | Show processes with the highest CPU load |
| `/help` | Show help message |
