package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
		SetAcceptanceFunc(tview.InputFieldMaxLength(50))

	renameInput.SetFieldBackgroundColor(tcell.ColorGray)

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
	//
	// modal := func(p tview.Primitive, width, height int) tview.Primitive {
	// 	return tview.NewGrid().
	// 		SetColumns(0, width, 0).
	// 		SetRows(0, height, 0).
	// 		AddItem(p, 1, 1, 1, 1, 0, 0, true)
	// }
	// Layout

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
				index := fileList.GetCurrentItem()
				fileName, _ := fileList.GetItemText(index)
				fileList.SetCurrentItem(index + 1)
				os.Remove(filepath.Join(directory, fileName))
				refreshFileList(fileList, directory)
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
