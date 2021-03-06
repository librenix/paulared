package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SummaryPage struct {
	*widgets.QWizardPage
	image *widgets.QLabel
	summary *widgets.QLabel
	desc *widgets.QLabel
}

func NewSummaryPage(parent widgets.QWidget_ITF) *SummaryPage {
	image := widgets.NewQLabel(nil, 0)
	summary := widgets.NewQLabel2("<h1>The installation was successful.</h1>", nil, 0)
	desc := widgets.NewQLabel2("The software was installed.", nil, 0)

	page := &SummaryPage{widgets.NewQWizardPage(parent), image, summary, desc}
	page.SetTitle("The installation was completed successfully.")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(image, 0, core.Qt__AlignBaseline)
	layout.AddWidget(summary, 0, core.Qt__AlignBaseline)
	layout.AddWidget(desc, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}