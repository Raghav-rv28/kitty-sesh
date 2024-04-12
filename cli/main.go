package main

import "github.com/rivo/tview"

func main() {
	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("List item 1", "Some explanatory text", '1', nil).
		AddItem("List item 2", "Some explanatory text", '2', nil).
		AddItem("List item 3", "Some explanatory text", '3', nil).
		AddItem("List item 4", "Some explanatory text", '4', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
