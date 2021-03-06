package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type DetectPage struct {
	*widgets.QWizardPage
	lv *widgets.QListView
}

func NewDetectPage(parent widgets.QWidget_ITF) *DetectPage {
	lv := widgets.NewQListView(nil)
	
	page := &DetectPage{widgets.NewQWizardPage(parent), lv}
	page.SetTitle("Detect hardware")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(lv, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}