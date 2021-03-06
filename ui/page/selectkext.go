package page

import (
	"github.com/therecipe/qt/core"
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

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(lv, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}