package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OSWindow struct {
	Tabs []Tab `json:"tabs"`
}

func main() {
	// check if the terminal is kitty
	term := os.Getenv("TERM")
	if term != "xterm-kitty" && term != "kitty" {
		fmt.Fprintf(os.Stderr, "error: this command must be run inside a kitty terminal\n")
		os.Exit(1)
	}
	// make the output file and folder.
	folderPath := getFolderPath()
	if len(os.Args) > 1 && os.Args[1] == "ss" {
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

		// Constant folderPath
		// get list gui with list of kitty sessions in the folderPath.
		fileList := getFileList(folderPath)
		// using frame to add instructions on top
		frame := tview.NewFrame(fileList).
			AddText("Press 'q' to Quit kitty-sesh; 'r' to Rename sessions", false, tview.AlignCenter, tcell.ColorWhite).
			AddText("'d' to Delete sessions; 'D' to delete All Sessions", false, tview.AlignCenter, tcell.ColorWhite).
			AddText("Use Arrow Keys to traverse the list", true, tview.AlignCenter, tcell.ColorWhite).
			AddText("Press Enter to start the session", true, tview.AlignCenter, tcell.ColorWhite)
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

		// grid := tview.NewGrid().SetRows(30, 0, 5).
		// 	AddItem(frame, 0, 0, 1, 1, 0, 125, true).
		// 	AddItem(fileContent, 1, 0, 1, 1, 0, 125, true)
		renameInput.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				oldName, _ := fileList.GetItemText(fileList.GetCurrentItem())
				newName := renameInput.GetText()
				newName += ".kitty"
				if newName != "" && newName != oldName {
					err := os.Rename(filepath.Join(folderPath, oldName), filepath.Join(folderPath, newName))
					if err != nil {
						fmt.Println("Error renaming file:", err)
					} else {
						refreshFileList(fileList, folderPath)
						fileContent.SetText("")
					}
				}
				flex.RemoveItem(renameInput)
				// grid.RemoveItem(renameInput)
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
					// grid.AddItem(renameInput, 2, 0, 1, 1, 0, 0, true)
					app.SetFocus(renameInput)
					// flex.AddItem(modal(renameInput, 75, 25), 1, 0, true)
					return nil
				case 'q', 'Q':
					app.Stop()
					return nil
				case 'd':
					modal.SetText("Are you sure you want to delete this session?")
					flex.AddItem(modal, 1, 0, true)
					app.SetFocus(modal)
					modal.SetDoneFunc(func(btnIndx int, btnLbl string) {
						if btnLbl == "Yes" {
							index := fileList.GetCurrentItem()
							fileName, _ := fileList.GetItemText(index)
							fileList.SetCurrentItem(index + 1)
							content, err := os.ReadFile(filepath.Join(folderPath, fileName))
							if err != nil {
								fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
							}
							fileContent.SetText(string(content))
							os.Remove(filepath.Join(folderPath, fileName))
							refreshFileList(fileList, folderPath)
						}
						flex.RemoveItem(modal)
						app.SetFocus(fileList)
					})
					return nil
				case 'D':
					modal.SetText("Are you sure you want to delete All sessions?")
					flex.AddItem(modal, 1, 0, true)
					app.SetFocus(modal)
					modal.SetDoneFunc(func(btnIndx int, btnLbl string) {
						if btnLbl == "Yes" {
							os.RemoveAll(folderPath)
							if err := os.MkdirAll(folderPath, 0755); err != nil {
								fmt.Fprintf(os.Stderr, "error creating folder: %v\n", err)
								os.Exit(1)
							}
							refreshFileList(fileList, folderPath)
						}
						flex.RemoveItem(modal)
						app.SetFocus(fileList)
					})

				}
				// updating the FileContent primitive to show the contents of the file.
			} else if event.Key() == tcell.KeyUp {
				if fileList.GetItemCount() == 0 {
					fileContent.SetText(fmt.Sprintln("No Sessions Available"))
					return nil
				}
				index := fileList.GetCurrentItem()
				if index > 0 {
					index -= 1
				} else if index == 0 {
					index = fileList.GetItemCount() - 1
				}
				fileName, _ := fileList.GetItemText(index)
				content, err := os.ReadFile(filepath.Join(folderPath, fileName))
				if err != nil {
					fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
					return nil
				}
				fileContent.SetText(string(content))
			} else if event.Key() == tcell.KeyDown {
				if fileList.GetItemCount() == 0 {
					fileContent.SetText(fmt.Sprintln("No Sessions Available"))
					return nil
				}
				index := fileList.GetCurrentItem()
				if index+1 < fileList.GetItemCount() {
					index += 1
				} else if index+1 == fileList.GetItemCount() {
					index = 0
				}
				fileName, _ := fileList.GetItemText(index)
				content, err := os.ReadFile(filepath.Join(folderPath, fileName))
				if err != nil {
					fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
					return nil
				}
				fileContent.SetText(string(content))
			}
			return event
		})

		// Load files from the folderPath
		refreshFileList(fileList, folderPath)

		// Set initial file content
		if fileList.GetItemCount() > 0 {
			firstFile, _ := fileList.GetItemText(0)
			fileContent.SetText(readFile(filepath.Join(folderPath, firstFile)))
		}

		if err := app.SetRoot(flex, true).Run(); err != nil {
			panic(err)
		}
	}
}
