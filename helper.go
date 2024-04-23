package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rivo/tview"
	"golang.org/x/crypto/ssh/terminal"
)

func getNextFileName(folderPath string, Args []string) (string, error) {
	// Get the list of files in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}
	var filename string
	if len(Args) > 2 {
		filename = fmt.Sprintf("%s.kitty", Args[2])
	} else {
		// Calculate the next file number
		fileNumber := len(files) + 1

		// Generate the filename with datetime and file number
		currentTime := time.Now().Format("2006-01-02T15-04-05")
		filename = fmt.Sprintf("%s-%d.kitty", currentTime, fileNumber)
	}
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
			case "stack":
				stackLayout(tab, outputFile)
			}
		}
	}
}

func refreshFileList(fileList *tview.List, folderPath string) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading folderPath:", err)
		return
	}

	fileList.Clear()

	for i, file := range files {
		fileList.AddItem(file.Name(), "", rune(i+49), nil)
	}
}

func readFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %s", err)
	}
	return string(content)
}

func getFileList(folderPath string, modalErr *tview.Modal, flex *tview.Flex) *tview.List {
	// List of files
	fileList := tview.NewList().
		ShowSecondaryText(false)
	fileList.SetTitle("Sessions").SetBorder(true)
	// Load files from the folderPath
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading folderPath:", err)
		return nil
	}

	for i, file := range files {
		fileList.AddItem(file.Name(), "", rune(i+49), nil)
	}

	// run the sessions on pressing enter
	fileList.SetSelectedFunc(func(index int, primaryString string, _ string, _ rune) {
		cmd := exec.Command("kitty", "--detach", "--session", filepath.Join(folderPath, primaryString))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			modalErr.SetText(fmt.Sprintln("Error launching Kitty session:", cmd.Args, err))
			modalErr.SetDoneFunc(func(btnIndex int, btnLabel string) {
				if btnLabel == "Close" {
					flex.RemoveItem(modalErr)
				}
			})
			flex.AddItem(modalErr, 1, 0, true)

			fmt.Println("Error launching Kitty session:", cmd.Args, err)
			return
		}
	})
	return fileList
}

func getFolderPath() string {
	usr, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home folderPath:", err)
		os.Exit(1)
	}

	// Define the folder path
	folderPath := filepath.Join(usr, ".config", "kitty", "sessions")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating folder: %v\n", err)
		os.Exit(1)
	}
	return folderPath
}

func getTerminalSize() (int, int, error) {
	i := int(os.Stdout.Fd())
	width, height, err := terminal.GetSize(i)
	fmt.Println(width, height)
	return width, height, err
}
