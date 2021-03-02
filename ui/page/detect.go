package page

import (
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

	return page
}