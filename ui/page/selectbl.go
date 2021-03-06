package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SelectBLPage struct {
	*widgets.QWizardPage
	ocb *widgets.QRadioButton
	cb *widgets.QRadioButton
	nb *widgets.QRadioButton
	gb *widgets.QGroupBox
}

func NewSelectBLPage(parent widgets.QWidget_ITF) *SelectBLPage {
	gb := widgets.NewQGroupBox(nil)
	ocb := widgets.NewQRadioButton2("OpenCore", gb)
	cb := widgets.NewQRadioButton2("Clover", gb)
	nb := widgets.NewQRadioButton2("None", gb)
	ocb.SetChecked(true)

	page := &SelectBLPage{widgets.NewQWizardPage(parent), ocb, cb, nb, gb}
	page.SetTitle("Select bootloader")

	gl := widgets.NewQVBoxLayout2(parent)
	gl.AddWidget(ocb, 0, core.Qt__AlignBaseline)
	gl.AddWidget(cb, 0, core.Qt__AlignBaseline)
	gl.AddWidget(nb, 0, core.Qt__AlignBaseline)
	gb.SetLayout(gl)
	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(gb, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}