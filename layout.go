package main

import (
	"fmt"
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
	MainBias       []float32          `json:"main_bias"`
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

func verticalLayout(tab Tab, outputFile *os.File) {
	windows := tab.Windows
	totalWindows := len(windows)
	traverseArr := getTraverseArr(tab)
	// creating windows
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)
		keyToCheck := strconv.Itoa(id)
		_, resize := tab.LayoutState.BiasedMap[keyToCheck]
		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		if 59/totalWindows > window.Rows && resize {
			outputFile.WriteString(fmt.Sprintf("resize_window shorter %d\n", (59/totalWindows)-window.Rows))
		} else if 59/totalWindows < window.Rows && resize {
			outputFile.WriteString(fmt.Sprintf("resize_window taller %d\n", window.Rows-(59/totalWindows)))
		}
		if window.IsFocused {
			fmt.Println("focus")
		}
	}
}

func horiztonalLayout(tab Tab, outputFile *os.File) {
	// grabbing the list of windows id in the order they are to be added
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	totalWindows := len(windows)
	for _, id := range traverseArr {
		window := getWindow(windows, id)
		if window == nil {
			continue
		}
		// make the command for each window
		cmd := loopBreak(window.Title)
		keyToCheck := strconv.Itoa(id)
		_, resize := tab.LayoutState.BiasedMap[keyToCheck]

		outputFile.WriteString(fmt.Sprintf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd))
		if 255/totalWindows > window.Cols && resize {
			outputFile.WriteString(fmt.Sprintf("resize_window narrower %d\n", (255/totalWindows)-window.Cols))
		} else if 255/totalWindows < window.Cols && resize {
			outputFile.WriteString(fmt.Sprintf("resize_window wider %d\n", window.Cols-(255/totalWindows)))
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
