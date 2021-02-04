package main

import (
	"fmt"
	"os"

	"github.com/librenix/paulared/ui"
)

func main() {
	// Logging
	f, err := LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer f.Close()

	ui.Setup()
}
