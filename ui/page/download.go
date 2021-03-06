package page

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type DownloadPage struct {
	*widgets.QWizardPage
	dpb *widgets.QProgressBar
	fpb *widgets.QProgressBar
}

func NewDownloadPage(parent widgets.QWidget_ITF) *DownloadPage {
	dpb := widgets.NewQProgressBar(nil)
	fpb := widgets.NewQProgressBar(nil)
	
	page := &DownloadPage{widgets.NewQWizardPage(parent), dpb, fpb}
	page.SetTitle("Download and write")
	page.SetCommitPage(true)

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(dpb, 0, core.Qt__AlignBaseline)
	layout.AddWidget(fpb, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	return page
}