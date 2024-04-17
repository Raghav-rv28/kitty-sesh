package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type OSWindow struct {
	Tabs []Tab `json:"tabs"`
}

func getNextFileName(folderPath string) (string, error) {
	// Get the list of files in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	// Calculate the next file number
	fileNumber := len(files) + 1

	// Generate the filename with datetime and file number
	currentTime := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("%s-%d.kitty", currentTime, fileNumber)
	fmt.Println(filename)
	// Return the full path to the file
	return filepath.Join(folderPath, filename), nil
}

func getEnvVars(env map[string]string) string {
	var sb strings.Builder
	for key, value := range env {
		sb.WriteString(fmt.Sprintf("--env %s=%s ", key, value))
	}
	return strings.TrimSpace(sb.String())
}

func getTraverseArr(tab Tab) []int {
	var temp []int
	for _, item := range tab.Groups {
		temp = append(temp, item.Windows...)
	}
	return temp
}

func getWindow(windows []Window, id int) *Window {
	for i := range windows {
		if windows[i].Id == id {
			return &windows[i]
		}
	}
	return nil
}

// breaking the side case where one of the windows have the command which creates a new session.
func loopBreak(title string) string {
	var cmd string
	if strings.Contains(title, "--session") {
		cmd = ""
	} else if title == "~" {
		cmd = ""
	} else {
		cmd = title
	}
	return cmd
}

func convert(session []OSWindow, outputFile *os.File) {
	first := true

	for _, osWindow := range session {
		if !first {
			fmt.Println("\nnew_os_window")
		} else {
			first = false
		}

		for _, tab := range osWindow.Tabs {
			outputFile.WriteString(fmt.Sprintf("new_tab %s\n", tab.Title))
			outputFile.WriteString(fmt.Sprintf("layout %s\n", tab.Layout))
			switch tab.Layout {
			case "horizontal":
				horiztonalLayout(tab, outputFile)
			case "vertical":
				verticalLayout(tab, outputFile)
			case "grid":
				gridLayout(tab, outputFile)
			case "split":
				splitLayout(tab, outputFile)
			case "tall":
				tallLayout(tab, outputFile)
			case "fat":
				fatLayout(tab, outputFile)
			default:
			}
		}
	}
}

func main() {
	// check if the terminal is kitty
	term := os.Getenv("TERM")
	if term != "xterm-kitty" && term != "kitty" {
		fmt.Fprintf(os.Stderr, "error: this command must be run inside a kitty terminal\n")
		os.Exit(1)
	}
	// get session data.
	cmd := exec.Command("kitty", "@", "ls")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running kitty command: %v\n", err)
		os.Exit(1)
	}
	var session []OSWindow
	if err := json.Unmarshal(output, &session); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding JSON: %v\n", err)
		os.Exit(1)
	}
	// make the output file and folder.
	folderPath := "/home/raghav/.config/kitty/sessions"
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating folder: %v\n", err)
		os.Exit(1)
	}

	// Generate the next filename
	filename, err := getNextFileName(folderPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating filename: %v\n", err)
		os.Exit(1)
	}
	// Open a new file for writing
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %v\n", err)
		os.Exit(1)
	}

	defer outputFile.Close()
	// convert json data to session file for kitty
	convert(session, outputFile)
}
