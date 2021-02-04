package ui

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

func Setup() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Hackintosh Installer")
	
	app.Exec()
}
