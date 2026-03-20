package services

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemStats struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	LoadAvg     *load.AvgStat
	Uptime      uint64
	Hostname    string
	NetworkIO   *net.IOCountersStat
}

func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, fmt.Errorf("tidak dapat membaca CPU usage")
}

func GetSystemStats() (*SystemStats, error) {
	cpuUsage, err := GetCPUUsage()
	if err != nil {
		return nil, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	netIO, err := net.IOCounters(false)
	var networkStats *net.IOCountersStat
	if err == nil && len(netIO) > 0 {
		networkStats = &netIO[0]
	}

	return &SystemStats{
		CPUUsage:    cpuUsage,
		MemoryUsage: memInfo.UsedPercent,
		DiskUsage:   diskInfo.UsedPercent,
		LoadAvg:     loadAvg,
		Uptime:      hostInfo.Uptime,
		Hostname:    hostInfo.Hostname,
		NetworkIO:   networkStats,
	}, nil
}

func FormatStats(stats *SystemStats) string {
	uptime := time.Duration(stats.Uptime) * time.Second

	msg := fmt.Sprintf(
		"📊 *System Statistics*\n\n"+
			"🖥 *Hostname:* `%s`\n"+
			"⏱ *Uptime:* `%s`\n\n"+
			"🔥 *CPU Usage:* `%.2f%%`\n"+
			"💾 *Memory Usage:* `%.2f%%`\n"+
			"💿 *Disk Usage:* `%.2f%%`\n\n"+
			"⚖ *Load Average:*\n"+
			"  • 1m: `%.2f`\n"+
			"  • 5m: `%.2f`\n"+
			"  • 15m: `%.2f`",
		stats.Hostname,
		uptime.String(),
		stats.CPUUsage,
		stats.MemoryUsage,
		stats.DiskUsage,
		stats.LoadAvg.Load1,
		stats.LoadAvg.Load5,
		stats.LoadAvg.Load15,
	)

	if stats.NetworkIO != nil {
		msg += fmt.Sprintf("\n\n🌐 *Network:*\n"+
			"  • Download: `%.2f MB`\n"+
			"  • Upload: `%.2f MB`",
			float64(stats.NetworkIO.BytesRecv)/1024/1024,
			float64(stats.NetworkIO.BytesSent)/1024/1024,
		)
	}

	return msg
}

func GetTopProcesses() (string, error) {
	cmd := "ps -eo pcpu,pmem,comm --sort=-pcpu --no-headers | head -n 5"
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return "⚠️ Tidak ada data proses.", nil
	}

	var sb strings.Builder
	sb.WriteString("```\n")
	sb.WriteString(fmt.Sprintf("%-6s %-6s %-15s\n", "CPU", "MEM", "COMMAND"))
	sb.WriteString(strings.Repeat("-", 30) + "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			cpu := fields[0] + "%"
			mem := fields[1] + "%"
			comm := strings.Join(fields[2:], " ")

			if len(comm) > 15 {
				comm = comm[:12] + "..."
			}
			sb.WriteString(fmt.Sprintf("%-6s %-6s %-15s\n", cpu, mem, comm))
		}
	}
	sb.WriteString("```")

	return sb.String(), nil
}