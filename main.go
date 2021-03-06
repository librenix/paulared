package main

import (
	"fmt"
	"os"

	"github.com/librenix/paulared/ui"
	"github.com/librenix/paulared/utils"
)

func main() {
	// Logging
	f, err := utils.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer f.Close()

	ui.NewApp()
}
