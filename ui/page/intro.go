package page

import (
	"github.com/therecipe/qt/widgets"
)

type IntroPage struct {
	*widgets.QWizardPage
}

func NewIntroPage(parent widgets.QWidget_ITF) *IntroPage {
	page := &IntroPage{widgets.NewQWizardPage(parent)}
	page.SetTitle("Welcome")
	page.SetSubTitle("Welcome to paulared hackintosh installer")

	return page
}
