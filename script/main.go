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
	MainBias       []float32 `json:"main_bias"`
	NumFullSizeWin int16     `json:"num_full_size_windows"`
	Pairs          *Pairs    `json:"pairs,omitempty"`
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

func envToStr(env map[string]string) string {
	var sb strings.Builder
	for key, value := range env {
		sb.WriteString(fmt.Sprintf("--env %s=%s ", key, value))
	}
	return strings.TrimSpace(sb.String())
}

func cmdlineToStr(cmdline []string) string {
	return strings.Join(cmdline, " ")
}

func verticalLayout(tab Tab) error {
	fmt.Printf("title %s", tab.Title)
	return nil
}

func horiztonalLayout(tab Tab) error {
	// grabbing the list of windows id in the order they are to be added
	var traverseArr []int
	for _, item := range tab.Groups {
		traverseArr = append(traverseArr, item.Windows...)
	}
	// creating windows
	windows := tab.Windows
	totalWindows := len(windows)
	for i, id := range traverseArr {
		var window *Window

		for i := range windows {
			if windows[i].Id == id {
				window = &windows[i]
				break
			}
		}
		fmt.Printf("launch %s --hold --stdin-source=@screen_scrollback --title '%s' --cwd %s %s\n", envToStr(window.Env), window.Title, window.Cwd, window.Title)
		if i%2 == 0 {
			if 255/totalWindows > window.Cols {
				fmt.Printf("resize_window narrower %d\n", (255/totalWindows)-window.Cols)
			} else {
				fmt.Printf("resize_window wider %d\n", window.Cols-(255/totalWindows))
			}
		}
		if window.IsFocused {
			fmt.Println("focus")
		}
	}
	return nil
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
			// if tab.Windows != nil {
			// 	fmt.Printf("cd %s\n", tab.Windows[0].Cwd)
			// }

			// SPLITS LAYOUT SIZING
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
			switch tab.Layout {
			case "horizontal":
				horiztonalLayout(tab)
			case "vertical":
				verticalLayout(tab)
			default:
			}
			// for i, w := range tab.Windows {
			// 	fmt.Printf("title %s\n", w.Title)
			// 	fmt.Printf("launch %s %s\n", envToStr(w.Env), fgProcToStr(w.ForegroundProcesses))
			// 	if i == 0 {
			// 		fmt.Println("resize_window wider 10")
			// 	}
			// 	if w.IsFocused {
			// 		fmt.Println("focus")
			// 	}
			// }
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
