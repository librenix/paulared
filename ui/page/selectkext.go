package page

import (
	"github.com/therecipe/qt/widgets"
)

type SelectKextPage struct {
	*widgets.QWizardPage
	lv *widgets.QListView
}

func NewSelectKextPage(parent widgets.QWidget_ITF) *SelectKextPage {
	lv := widgets.NewQListView(nil)
	
	page := &SelectKextPage{widgets.NewQWizardPage(parent), lv}
	page.SetTitle("Select kexts")

	return page
}