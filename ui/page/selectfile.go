package page

import "github.com/therecipe/qt/widgets"

type SelectFilePage struct {
	*widgets.QWizardPage
}

func NewSelectFilePage(parent widgets.QWidget_ITF) *SelectFilePage {
	page := &SelectFilePage{widgets.NewQWizardPage(parent)}

	return page
}
