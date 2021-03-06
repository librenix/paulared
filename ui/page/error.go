package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type ErrorPage struct {
	*widgets.QWizardPage
	image *widgets.QLabel
	summary *widgets.QLabel
	desc *widgets.QLabel
}

func NewErrorPage(parent widgets.QWidget_ITF) *ErrorPage {
	image := widgets.NewQLabel(nil, 0)
	summary := widgets.NewQLabel2("<h1>The installation failed.</h1>", nil, 0)
	desc := widgets.NewQLabel2("The installer encountered an error that caused the installation to fail, Contact the software manufacturer for assistance.", nil, 0)

	page := &ErrorPage{widgets.NewQWizardPage(parent), image, summary, desc}
	page.SetTitle("An error occurred during installation.")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(image, 0, core.Qt__AlignBaseline)
	layout.AddWidget(summary, 0, core.Qt__AlignBaseline)
	layout.AddWidget(desc, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}