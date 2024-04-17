package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

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
	fileContent.SetBorder(true)
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
	fileList.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
		filename, _ := fileList.GetItemText(index)
		content, err := os.ReadFile(filepath.Join(directory, filename))
		if err != nil {
			fileContent.SetText(fmt.Sprintf("Error reading file: %s", err))
			return
		}
		fileContent.SetText(string(content))
	})

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(fileList, 0, 1, true).
		AddItem(fileContent, 0, 3, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
