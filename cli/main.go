package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	for _, file := range files {
		fileList.AddItem(file.Name(), "", 0, nil)
	}
}

func readFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %s", err)
	}
	return string(content)
}

func main() {
	app := tview.NewApplication()

	// Constant directory
	directory := "/home/raghav/.config/kitty/sessions/"

	// List of files
	fileList := tview.NewList().
		ShowSecondaryText(false)
	fileList.SetTitle("Sessions").SetBorder(true)

	// File content preview
	fileContent := tview.NewTextView()
	fileContent.SetBorder(true).SetTitle("Preview")
	// Load files from the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for i, file := range files {
		fileList.AddItem(file.Name(), "", rune(i+49), nil)
	}

	// Set event handler for file selection
	fileList.SetSelectedFunc(func(index int, primaryString string, _ string, _ rune) {
		content, err := os.ReadFile(filepath.Join(directory, primaryString))
		if err != nil {
			fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
			return
		}
		fileContent.SetText(string(content))
	})

	renameInput := tview.NewInputField().
		SetLabel("New Name: ").
		SetFieldWidth(60).
		SetAcceptanceFunc(tview.InputFieldMaxLength(50))

	renameInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			oldName, _ := fileList.GetItemText(fileList.GetCurrentItem())
			newName := renameInput.GetText()

			if newName != "" && newName != oldName {
				err := os.Rename(filepath.Join(directory, oldName), filepath.Join(directory, newName))
				if err != nil {
					fmt.Println("Error renaming file:", err)
				} else {
					refreshFileList(fileList, directory)
					fileContent.SetText("")
				}
			}

			app.SetFocus(fileList)
			fileList.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
				fileName, _ := fileList.GetItemText(index)
				fileContent.SetText(readFile(filepath.Join(directory, fileName)))
				renameInput.SetText(fileName)
			})
		}
	})
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}
	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(fileList, 0, 1, true).
		AddItem(fileContent, 0, 3, false)

	// Set Input Capture to handle custom key events
	fileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'r', 'R':
				index := fileList.GetCurrentItem()
				fileName, _ := fileList.GetItemText(index)
				renameInput.SetText(fileName)
				flex.AddItem(modal(renameInput, 75, 25), 1, 0, false)
				return nil
			}
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
