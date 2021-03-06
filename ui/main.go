package ui

import (
	"os"
	"runtime"

	"github.com/therecipe/qt/widgets"
)

type window struct {
	*widgets.QApplication
}

func NewApp() {
	// Application
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("paulared")
	app.SetApplicationName("paulared")

	// Main window
	window := widgets.NewQMainWindow(nil, 0)

	// Wizard
	wizard := NewWizard(nil, 0)
	window.SetCentralWidget(wizard)
	window.SetMinimumSize2(800, 640)
	window.SetWindowTitle("Hackintosh Installer")

	// Style
	switch (runtime.GOOS) {
	case "darwin":
		app.SetStyle2("macintosh")
		break
	case "windows":
		app.SetStyle2("windowsvista")
		break
	default:
		app.SetStyle2("fusion")
		break
	}

	window.Show()

	app.Exec()
}
