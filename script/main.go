package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ForegroundProcesses struct {
	Pid     int32    `json:"pid"`
	Cmdline []string `json:"cmdline"`
	Cwd     string   `json:"cwd"`
}

type Window struct {
	Title               string                `json:"title"`
	Env                 map[string]string     `json:"env"`
	ForegroundProcesses []ForegroundProcesses `json:"foreground_processes"`
	IsFocused           bool                  `json:"is_focused"`
	Cwd                 string                `json:"cwd"`
}

type Tab struct {
	Title   string   `json:"title"`
	Windows []Window `json:"windows"`
	Layout  string   `json:"layout"`
}

type OSWindow struct {
	Tabs []Tab `json:"tabs"`
}

func envToStr(env map[string]string) string {
	var sb strings.Builder
	for key, value := range env {
		sb.WriteString(fmt.Sprintf("--env %s=%s ", key, value))
	}
	return strings.TrimSpace(sb.String())
}

func cmdlineToStr(cmdline []string) string {
	return strings.Join(cmdline, " ")
}

func fgProcToStr(fg []ForegroundProcesses) string {
	var commands string
	for _, process := range fg {
		commands += cmdlineToStr(process.Cmdline)

		if strings.Contains(commands, "kitty @ ls") {
			return os.Getenv("SHELL")
		}
	}
	return strings.TrimSpace(commands)
}

func convert(session []OSWindow) {
	first := true

	for _, osWindow := range session {
		if !first {
			fmt.Println("\nnew_os_window")
		} else {
			first = false
		}

		for _, tab := range osWindow.Tabs {
			fmt.Printf("new_tab %s\n", tab.Title)
			fmt.Printf("layout %s\n", tab.Layout)
			if tab.Windows != nil {
				fmt.Printf("cd %s\n", tab.Windows[0].Cwd)
			}
			for _, w := range tab.Windows {
				fmt.Printf("title %s\n", w.Title)
				fmt.Printf("launch %s %s\n", envToStr(w.Env), fgProcToStr(w.ForegroundProcesses))
				if w.IsFocused {
					fmt.Println("focus")
				}
			}
		}
	}
}

func main() {
	var session []OSWindow
	if err := json.NewDecoder(os.Stdin).Decode(&session); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	convert(session)
}
