package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OSWindow struct {
	Tabs []Tab `json:"tabs"`
}

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

func refreshFileList(fileList *tview.List, directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
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

func getFileList(directory string) *tview.List {
	// List of files
	fileList := tview.NewList().
		ShowSecondaryText(false)
	fileList.SetTitle("Sessions").SetBorder(true)
	// Load files from the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	for i, file := range files {
		fileList.AddItem(file.Name(), "", rune(i+49), nil)
	}

	// run the sessions on pressing enter
	fileList.SetSelectedFunc(func(index int, primaryString string, _ string, _ rune) {
		cmd := exec.Command("kitty", "--detach", "--session", filepath.Join(directory, primaryString))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error launching Kitty session:", cmd.Args, err)
			return
		}
	})
	return fileList
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "ss" {
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
		usr, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			return
		}

		// Define the folder path
		folderPath := filepath.Join(usr, ".config", "kitty", "sessions")
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating folder: %v\n", err)
			os.Exit(1)
		}
		// Generate the next filename
		filename, err := getNextFileName(folderPath, os.Args)
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
	} else {
		app := tview.NewApplication()

		// Constant directory
		directory := "/home/raghav/.config/kitty/sessions/"
		// get list gui with list of kitty sessions in the directory.
		fileList := getFileList(directory)
		// using frame to add instructions on top
		frame := tview.NewFrame(fileList).
			AddText("Press 'q' to Quit kitty-sesh; 'r' to Rename sessions; 'd' to Delete sessions", false, tview.AlignCenter, tcell.ColorWhite).
			AddText("Use Arrow Keys to traverse the list", true, tview.AlignCenter, tcell.ColorWhite).
			AddText("Press Enter to start the session ", true, tview.AlignCenter, tcell.ColorWhite)
		// File content preview
		fileContent := tview.NewTextView()
		fileContent.SetBorder(true).SetTitle("Preview")

		renameInput := tview.NewInputField().
			SetLabel("New Name: ").
			SetFieldWidth(60).
			SetAcceptanceFunc(tview.InputFieldMaxLength(50)).
			SetFieldBackgroundColor(tcell.ColorGray)

		// modal to confirm changes
		modal := tview.NewModal().AddButtons([]string{"Yes", "No"})
		// the layout master
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(frame, 0, 1, true).
			AddItem(fileContent, 0, 3, false)

		renameInput.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				oldName, _ := fileList.GetItemText(fileList.GetCurrentItem())
				newName := renameInput.GetText()
				newName += ".kitty"
				if newName != "" && newName != oldName {
					err := os.Rename(filepath.Join(directory, oldName), filepath.Join(directory, newName))
					if err != nil {
						fmt.Println("Error renaming file:", err)
					} else {
						refreshFileList(fileList, directory)
						fileContent.SetText("")
					}
				}
				flex.RemoveItem(renameInput)
				app.SetFocus(fileList)

			}
		})

		// Set Input Capture to handle custom key events
		fileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'r', 'R':
					index := fileList.GetCurrentItem()
					fileName, _ := fileList.GetItemText(index)
					renameInput.SetText(strings.Split(fileName, ".kitty")[0])
					flex.AddItem(renameInput, 0, 1, true)
					app.SetFocus(renameInput)
					// flex.AddItem(modal(renameInput, 75, 25), 1, 0, true)
					return nil
				case 'q', 'Q':
					app.Stop()
					return nil
				case 'd', 'D':
					modal.SetText("Are you sure you want to delete this session?")
					flex.AddItem(modal, 1, 0, true)
					app.SetFocus(modal)
					modal.SetDoneFunc(func(btnIndx int, btnLbl string) {
						if btnLbl == "Yes" {
							index := fileList.GetCurrentItem()
							fileName, _ := fileList.GetItemText(index)
							fileList.SetCurrentItem(index + 1)
							content, err := os.ReadFile(filepath.Join(directory, fileName))
							if err != nil {
								fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
							}
							fileContent.SetText(string(content))
							os.Remove(filepath.Join(directory, fileName))
							refreshFileList(fileList, directory)
						}
						flex.RemoveItem(modal)
						app.SetFocus(fileList)
					})
					return nil
				}
				// updating the FileContent primitive to show the contents of the file.
			} else if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
				if fileList.GetItemCount() == 0 {
					fileContent.SetText(fmt.Sprintln("No Sessions Available"))
					return nil
				}
				index := fileList.GetCurrentItem()
				fileName, _ := fileList.GetItemText(index)
				content, err := os.ReadFile(filepath.Join(directory, fileName))
				if err != nil {
					fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
					return nil
				}
				fileContent.SetText(string(content))
			}
			return event
		})

		// Load files from the directory
		refreshFileList(fileList, directory)

		// Set initial file content
		if fileList.GetItemCount() > 0 {
			firstFile, _ := fileList.GetItemText(0)
			fileContent.SetText(readFile(filepath.Join(directory, firstFile)))
		}

		if err := app.SetRoot(flex, true).Run(); err != nil {
			panic(err)
		}
	}
}
