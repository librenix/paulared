package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SelectKextPage struct {
	*widgets.QWizardPage
	lw *widgets.QListWidget
}

func NewSelectKextPage(parent widgets.QWidget_ITF) *SelectKextPage {
	lw := widgets.NewQListWidget(nil)

	page := &SelectKextPage{widgets.NewQWizardPage(parent), lw}
	page.SetTitle("Select kexts")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(lw, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}
