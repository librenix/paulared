package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mattn/go-isatty"
)

var version = "0.1.0"

func main() {
	// Logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer f.Close()

	// Is cli or gui?
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		// Execute cli commands
		if len(os.Args) > 1 {
			Execute()
		} else {
			// Start TUI if without any args
			fmt.Println(os.Args)
			p := SetupUI()
			if err := p.Start(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		// Qt GUI
		//qt.Setup()
	}
}
