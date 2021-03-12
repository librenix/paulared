package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SelectFilePage struct {
	*widgets.QWizardPage
	gb      *widgets.QGroupBox
	online  *widgets.QRadioButton
	offline *widgets.QRadioButton
	oslist  *widgets.QListWidget
	path    *widgets.QLineEdit
	browse  *widgets.QPushButton
}

func NewSelectFilePage(parent widgets.QWidget_ITF) *SelectFilePage {
	macver := []string{
		"macOS 11 Big Sur (latest)",
		"macOS 10.15 Catalina",
		"macOS 10.14 Mojav",
		"macOS 10.13 High Sierra",
		"macOS 10.12 Sierra",
		"OS X 10.11 El Capitan",
		"OS X 10.10 Yosemite",
		"OS X 10.9 Mavericks",
		"OS X 10.8 Mountain Lion",
		"OS X 10.7 Lion",
		// Mac OS X 10.6 Snow Leopard,
		// Mac OS X 10.5 Leopard,
		// Mac OS X 10.4 Tiger,
	}

	gb := widgets.NewQGroupBox(nil)
	online := widgets.NewQRadioButton2("Online installation", nil)
	offline := widgets.NewQRadioButton2("Offline installation", nil)
	oslist := widgets.NewQListWidget(nil)
	path := widgets.NewQLineEdit(nil)
	browse := widgets.NewQPushButton2("Browse...", nil)
	oslist.AddItems(macver)
	online.SetChecked(true)

	page := &SelectFilePage{widgets.NewQWizardPage(parent), gb, online, offline, oslist, path, browse}
	page.SetTitle("Select dmg file")

	fl := widgets.NewQHBoxLayout2(parent)
	fl.AddWidget(path, 0, core.Qt__AlignBaseline)
	fl.AddWidget(browse, 0, core.Qt__AlignBaseline)
	file := widgets.NewQFrame(nil, 0)
	file.SetLayout(fl)
	gl := widgets.NewQVBoxLayout2(parent)
	gl.AddWidget(online, 0, core.Qt__AlignBaseline)
	gl.AddWidget(oslist, 0, core.Qt__AlignBaseline)
	gl.AddWidget(offline, 0, core.Qt__AlignBaseline)
	gl.AddWidget(file, 0, core.Qt__AlignBaseline)
	gb.SetLayout(gl)
	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(gb, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}
