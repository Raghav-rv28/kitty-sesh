package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// Structs
type Pairs struct {
	Bias       float32     `json:"bias"`
	Horizontal bool        `json:"horizontal"`
	One        interface{} `json:"one"`
	Two        interface{} `json:"two"`
}

type LayoutState struct {
	MainBias       []float64          `json:"main_bias"`
	NumFullSizeWin int16              `json:"num_full_size_windows"`
	Pairs          *Pairs             `json:"pairs,omitempty"`
	BiasedCols     map[string]float64 `json:"biased_cols"`
	BiasedRows     map[string]float64 `json:"biased_rows"`
	BiasedMap      map[string]float64 `json:"biased_map"`
}

type Window struct {
	Title     string            `json:"title"`
	Env       map[string]string `json:"env"`
	IsFocused bool              `json:"is_focused"`
	Cwd       string            `json:"cwd"`
	Rows      int               `json:"lines"`
	Cols      int               `json:"columns"`
	Id        int               `json:"id"`
}

type Groups struct {
	Id      int   `json:"id"`
	Windows []int `json:"windows"`
}
type Tab struct {
	Title       string      `json:"title"`
	Windows     []Window    `json:"windows"`
	Layout      string      `json:"layout"`
	LayoutState LayoutState `json:"layout_state"`
	Groups      []Groups    `json:"groups"`
}

func resizeLayoutBased(tab Tab, outputFile *os.File, window *Window) {
	rows, columns, err := getTerminalSize()
	fmt.Println("rows:", rows, "columns:", columns)
	if err != nil {
		panic(err)
	}
	layout := tab.Layout
	totalWindows := len(tab.Windows)
	switch layout {
	case "horizontal":
		if columns/totalWindows > window.Cols {
			outputFile.WriteString(fmt.Sprintf("resize_window narrower %d\n", (columns/totalWindows)-window.Cols))
		} else if columns/totalWindows < window.Cols {
			outputFile.WriteString(fmt.Sprintf("resize_window wider %d\n", window.Cols-(columns/totalWindows)))
		}
	case "vertical":
		if rows/totalWindows > window.Rows {
			outputFile.WriteString(fmt.Sprintf("resize_window shorter %d\n", (rows/totalWindows)-window.Rows))
		} else if rows/totalWindows < window.Rows {
			outputFile.WriteString(fmt.Sprintf("resize_window taller %d\n", window.Rows-(rows/totalWindows)))
		}
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

func verticalLayout(tab Tab, outputFile *os.File) {
	windows := tab.Windows
	traverseArr := getTraverseArr(tab)
	// creating windows
	for indx, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)
		keyToCheck := strconv.Itoa(indx)
		_, resize := tab.LayoutState.BiasedMap[keyToCheck]
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		if resize {
			resizeLayoutBased(tab, outputFile, window)
		}
		if window.IsFocused {
			outputFile.WriteString(fmt.Sprintln("focus"))
		}
	}
}

func horiztonalLayout(tab Tab, outputFile *os.File) {
	// grabbing the list of windows id in the order they are to be added
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	for indx, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)
		keyToCheck := strconv.Itoa(indx)
		_, resize := tab.LayoutState.BiasedMap[keyToCheck]

		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		if resize {
			resizeLayoutBased(tab, outputFile, window)
		}
		if window.IsFocused {
			fmt.Println("focus")
		}
	}
}

func gridLayout(tab Tab, outputFile *os.File) {
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window

		cmd := loopBreak(window.Title)

		// Resizing stuff
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		// keys := make([]int, 0, totalWindows)
		// for k := range tab.LayoutState.BiasedCols {
		// 	intKey, err := strconv.Atoi(k)
		// 	if err != nil {
		// 		fmt.Println("Error converting key to int:", err)
		// 		return
		// 	}
		// 	keys = append(keys, intKey)
		// }
		// if indx == keys[indx] || (indx != 0 && len(keys) > indx/2 && keys[indx/2] != 0) {
		//   fmt.Println("resize_window ")
		//   }
	}
}

func splitLayout(tab Tab, outputFile *os.File) {
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window

		cmd := loopBreak(window.Title)

		// Resizing stuff
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		//SPLITS LAYOUT SIZING
		// if tab.LayoutState.Pairs != nil {
		// 	outputFile.WriteString(fmt.Sprintf("bias %f\n", tab.LayoutState.Pairs.Bias)
		// 	outputFile.WriteString(fmt.Sprintf("horizontal %t\n", tab.LayoutState.Pairs.Horizontal)
		//
		// 	// Print type and value of One
		// 	if intValue, ok := tab.LayoutState.Pairs.One.(float64); ok {
		// 		outputFile.WriteString(fmt.Sprintf("pair one: %f\n", intValue)
		// 	} else if pairs, ok := tab.LayoutState.Pairs.One.(map[string]interface{}); ok {
		// 		if one, ok := pairs["one"].(int); ok {
		// 			outputFile.WriteString(fmt.Sprintf("pair one: %d\n", one)
		// 		}
		// 	}
		// 	// Print type and value of Two
		// 	if intValue, ok := tab.LayoutState.Pairs.Two.(float64); ok {
		// 		outputFile.WriteString(fmt.Sprintf("pair two: %f\n", intValue)
		// 	} else if pairs, ok := tab.LayoutState.Pairs.Two.(map[string]interface{}); ok {
		// 		outputFile.WriteString(fmt.Sprintf("pair two: %.2f\n", pairs["bias"].(float64))
		// 	}
		// }
	}
}

func tallLayout(tab Tab, outputFile *os.File) {
	bias := tab.LayoutState.MainBias[0]
	for i := 1; i < int(tab.LayoutState.NumFullSizeWin); i++ {
		bias += tab.LayoutState.MainBias[i]
	}

	outputFile.WriteString(fmt.Sprintf("layout tall:bias=%d;full_size=%d\n", int(math.Round(bias*100)), tab.LayoutState.NumFullSizeWin))
	traverseArr := getTraverseArr(tab)
	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)

		// Resizing stuff
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
	}
}

func fatLayout(tab Tab, outputFile *os.File) {
	traverseArr := getTraverseArr(tab)
	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)

		// Resizing stuff
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
	}
}

func stackLayout(tab Tab, outputFile *os.File) {
	traverseArr := getTraverseArr(tab)
	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)

		// Resizing stuff
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
	}
}
