package ui

import (
	"github.com/librenix/paulared/ui/page"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Wizard struct {
	*widgets.QWizard
}

func NewWizard(parent widgets.QWidget_ITF, ff core.Qt__WindowType) *Wizard {
	wizard := &Wizard{widgets.NewQWizard(parent, ff)}

	// Create and add wizard pages
	intro := page.NewIntroPage(wizard)
	detect := page.NewDetectPage(wizard)
	selectdisk := page.NewSelectDiskPage(wizard)
	selectbl := page.NewSelectBLPage(wizard)
	selectkext := page.NewSelectKextPage(wizard)
	selectfile := page.NewSelectFilePage(wizard)
	download := page.NewDownloadPage(wizard)
	failed := page.NewErrorPage((wizard))
	summary := page.NewSummaryPage(wizard)

	wizard.AddPage(intro)
	wizard.AddPage(detect)
	wizard.AddPage(selectdisk)
	wizard.AddPage(selectbl)
	wizard.AddPage(selectkext)
	wizard.AddPage(selectfile)
	wizard.AddPage(download)
	wizard.AddPage(failed)
	wizard.AddPage(summary)

	wizard.SetOptions(widgets.QWizard__DisabledBackButtonOnLastPage | widgets.QWizard__NoCancelButton)
	wizard.SetWizardStyle(widgets.QWizard__MacStyle)

	return wizard
}
