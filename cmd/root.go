/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Monit struct {
	Memory int     `json:"memory"`
	Cpu    float64 `json:"cpu"`
}

type Process struct {
	Pid     int                    `json:"pid"`
	Name    string                 `json:"name"`
	Pm_id   int                    `json:"pm_id"`
	Monit   Monit                  `json:"monit"`
	Pm2_env map[string]interface{} `json:"pm2_env"`
}

func Execute() {
	out, err := exec.Command("pm2", "jlist").Output()
	if err != nil {
		fmt.Println(err)
	}

	var processes []Process

	json.Unmarshal(out, &processes)

	header := color.New(color.FgHiCyan, color.Bold)
	online := color.New(color.FgGreen, color.Bold)
	offline := color.New(color.FgHiRed, color.Bold)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{header.Sprint("id"), header.Sprint("name"), header.Sprint("pid"), header.Sprint("uptime"), header.Sprint("↺"), header.Sprint("status"), header.Sprint("cpu"), header.Sprint("memory")})
	t.SetStyle(table.Style{
		Name:    "pm2",
		Box:     table.StyleBoxRounded,
		Options: table.OptionsDefault,
		Format: table.FormatOptions{
			Header: text.Format(text.AlignLeft),
		},
	})

	if len(processes) == 0 {
		t.Render()
		return
	}

	for _, process := range processes {

		pmUptime := time.Unix(0, int64(process.Pm2_env["pm_uptime"].(float64))*int64(time.Millisecond))

		var pm_uptime string
		if process.Pm2_env["status"] == "online" {
			pm_uptime = formatDuration(time.Since(pmUptime))
		} else {
			pm_uptime = "0"
		}

		var pm_status string
		if process.Pm2_env["status"] == "online" {
			pm_status = online.Sprint("online")
		} else {
			pm_status = offline.Sprint("stopped")
		}

		t.AppendRow(table.Row{header.Sprint(process.Pm_id), process.Name, process.Pid, pm_uptime, process.Pm2_env["restart_time"], pm_status, fmt.Sprint(fmt.Sprintf("%.0f", process.Monit.Cpu), "%"), formatMemory(process.Monit.Memory)})
	}

	t.Render()
}

func formatMemory(memory int) string {
	if memory > 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(memory)/1024/1024/1024)
	} else if memory > 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(memory)/1024/1024)
	} else if memory > 1024 {
		return fmt.Sprintf("%.2f KB", float64(memory)/1024)
	} else {
		return fmt.Sprintf("%d B", memory)
	}
}

func formatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dD", days)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}
