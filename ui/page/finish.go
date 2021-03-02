package page

import (
	"github.com/therecipe/qt/widgets"
)

type FinishPage struct {
	*widgets.QWizardPage
	text *widgets.QLabel
}

func NewFinishPage(parent widgets.QWidget_ITF) *FinishPage {
	text := widgets.NewQLabel2("", nil, 0)

	page := &FinishPage{widgets.NewQWizardPage(parent), text}
	page.SetTitle("Finish")

	return page
}