package page

import (
	"github.com/therecipe/qt/widgets"
)

type SelectBLPage struct {
	*widgets.QWizardPage
	ocb *widgets.QCheckBox
	cb *widgets.QCheckBox
	nb *widgets.QCheckBox
	gb *widgets.QGroupBox
}

func NewSelectBLPage(parent widgets.QWidget_ITF) *SelectBLPage {
	gb := widgets.NewQGroupBox(nil)
	ocb := widgets.NewQCheckBox2("OpenCore", gb)
	cb := widgets.NewQCheckBox2("Clover", gb)
	nb := widgets.NewQCheckBox2("None", gb)

	page := &SelectBLPage{widgets.NewQWizardPage(parent), ocb, cb, nb, gb}
	page.SetTitle("Select bootloader")

	return page
}