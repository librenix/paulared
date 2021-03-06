package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SelectDiskPage struct {
	*widgets.QWizardPage
	lv *widgets.QListView
}

func NewSelectDiskPage(parent widgets.QWidget_ITF) *SelectDiskPage {
	lv := widgets.NewQListView(nil)

	page := &SelectDiskPage{widgets.NewQWizardPage(parent), lv}
	page.SetTitle("Select disk")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(lv, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}