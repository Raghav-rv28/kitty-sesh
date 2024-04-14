package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

type OSWindow struct {
	Tabs []Tab `json:"tabs"`
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

func getWindow(windows []Window, window *Window, id int) {
	for i := range windows {
		if windows[i].Id == id {
			window = &windows[i]
			break
		}
	}
}

// breaking the side case where one of the windows have the command which creates a new session.
func loopBreak(title string) string {
	var cmd string
	if strings.Contains(title, "--session") {
		cmd = ""
	} else {
		cmd = title
	}
	return cmd
}

func verticalLayout(tab Tab) {
	windows := tab.Windows
	totalWindows := len(windows)
	traverseArr := getTraverseArr(tab)
	// creating windows
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)
		// make the command for each window
		cmd := loopBreak(window.Title)
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
		if 59/totalWindows > window.Rows {
			fmt.Printf("resize_window shorter %d\n", (59/totalWindows)-window.Rows)
		} else if 59/totalWindows < window.Rows {
			fmt.Printf("resize_window taller %d\n", window.Rows-(59/totalWindows))
		}
		if window.IsFocused {
			fmt.Println("focus")
		}
	}
}

func horiztonalLayout(tab Tab) {
	// grabbing the list of windows id in the order they are to be added
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	totalWindows := len(windows)
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)
		// make the command for each window
		cmd := loopBreak(window.Title)

		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
		if 255/totalWindows > window.Cols {
			fmt.Printf("resize_window narrower %d\n", (255/totalWindows)-window.Cols)
		} else if 255/totalWindows < window.Cols {
			fmt.Printf("resize_window wider %d\n", window.Cols-(255/totalWindows))
		}
		if window.IsFocused {
			fmt.Println("focus")
		}
	}
}

func gridLayout(tab Tab) {
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)
		// make the command for each window

		cmd := loopBreak(window.Title)

		// Resizing stuff
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
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

func splitLayout(tab Tab) {
	traverseArr := getTraverseArr(tab)

	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)

		// make the command for each window

		cmd := loopBreak(window.Title)

		// Resizing stuff
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
		//SPLITS LAYOUT SIZING
		// if tab.LayoutState.Pairs != nil {
		// 	fmt.Printf("bias %f\n", tab.LayoutState.Pairs.Bias)
		// 	fmt.Printf("horizontal %t\n", tab.LayoutState.Pairs.Horizontal)
		//
		// 	// Print type and value of One
		// 	if intValue, ok := tab.LayoutState.Pairs.One.(float64); ok {
		// 		fmt.Printf("pair one: %f\n", intValue)
		// 	} else if pairs, ok := tab.LayoutState.Pairs.One.(map[string]interface{}); ok {
		// 		if one, ok := pairs["one"].(int); ok {
		// 			fmt.Printf("pair one: %d\n", one)
		// 		}
		// 	}
		// 	// Print type and value of Two
		// 	if intValue, ok := tab.LayoutState.Pairs.Two.(float64); ok {
		// 		fmt.Printf("pair two: %f\n", intValue)
		// 	} else if pairs, ok := tab.LayoutState.Pairs.Two.(map[string]interface{}); ok {
		// 		fmt.Printf("pair two: %.2f\n", pairs["bias"].(float64))
		// 	}
		// }
	}
}

func tallLayout(tab Tab) {
	traverseArr := getTraverseArr(tab)
	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)
		// make the command for each window
		cmd := loopBreak(window.Title)

		// Resizing stuff
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
	}
}

func fatLayout(tab Tab) {
	traverseArr := getTraverseArr(tab)
	// creating windows
	windows := tab.Windows
	// totalWindows := len(windows)
	for _, id := range traverseArr {
		var window *Window
		getWindow(windows, window, id)
		// make the command for each window
		cmd := loopBreak(window.Title)

		// Resizing stuff
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", getEnvVars(window.Env), window.Title, window.Cwd, cmd)
	}
}

func convert(session []OSWindow) {
	first := true

	for _, osWindow := range session {
		if !first {
			fmt.Println("\nnew_os_window")
		} else {
			first = false
		}

		for _, tab := range osWindow.Tabs {
			fmt.Printf("new_tab %s\n", tab.Title)
			fmt.Printf("layout %s\n", tab.Layout)
			switch tab.Layout {
			case "horizontal":
				horiztonalLayout(tab)
			case "vertical":
				verticalLayout(tab)
			case "grid":
				gridLayout(tab)
			case "split":
				splitLayout(tab)
			case "tall":
				tallLayout(tab)
			case "fat":
				fatLayout(tab)
			default:
			}
		}
	}
}

func main() {
	var session []OSWindow
	if err := json.NewDecoder(os.Stdin).Decode(&session); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	convert(session)
}
