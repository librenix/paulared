package page

import (
	"github.com/therecipe/qt/widgets"
)

type DownloadPage struct {
	*widgets.QWizardPage
	pb *widgets.QProgressBar
}

func NewDownloadPage(parent widgets.QWidget_ITF) *DownloadPage {
	pb := widgets.NewQProgressBar(nil)
	
	page := &DownloadPage{widgets.NewQWizardPage(parent), pb}
	page.SetTitle("Download and write")

	return page
}