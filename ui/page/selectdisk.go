package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SelectDiskPage struct {
	*widgets.QWizardPage
	lw *widgets.QListWidget
}

func NewSelectDiskPage(parent widgets.QWidget_ITF) *SelectDiskPage {
	lw := widgets.NewQListWidget(nil)
	lw.SetViewMode(widgets.QListView__IconMode)

	page := &SelectDiskPage{widgets.NewQWizardPage(parent), lw}
	page.SetTitle("Select disk")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(lw, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}
