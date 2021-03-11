package page

import (
	"github.com/google/gousb"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"

	c "github.com/librenix/paulared/core"
)

type DetectPage struct {
	*widgets.QWizardPage
	tw *widgets.QTreeWidget
}

func NewDetectPage(parent widgets.QWidget_ITF) *DetectPage {
	tw := widgets.NewQTreeWidget(nil)
	tw.SetColumnCount(3)

	page := &DetectPage{widgets.NewQWizardPage(parent), tw}
	page.SetTitle("Detect hardware")

	layout := widgets.NewQVBoxLayout2(parent)
	layout.AddWidget(tw, 0, core.Qt__AlignBaseline)
	page.SetLayout(layout)

	page.Init()

	return page
}

func (p *DetectPage) Init() {
	p.tw.SetHeaderLabels([]string{"type", "model", "status"})

	// Get hardware info
	hw := c.Detect()
	ctx := *gousb.NewContext()
	host := map[string]interface{}{
		"CPU":         *hw.CPU,
		"Methodboard": *hw.Product,
		"Memory":      *hw.Memory,
		"Disk":        *hw.Block,
		"GPU":         *hw.GPU,
		"PCI":         *hw.PCI,
		"USB":         ctx,
	}

	// Create a slice to order keys
	keys := make([]string, 0, len(host))
	for k, _ := range host {
		keys = append(keys, k)
	}

	// Add info to tree
	for k, v := range keys {
		res := c.Compatibility(host[v])
		model := res.Model
		var status string
		b := gui.NewQBrush()

		switch res.Status {
		case c.Supported:
			status = "Supported"
			b.SetColor(gui.NewQColor2(core.Qt__green))
			break
		case c.Unsupported:
			status = "Unsupported"
			b.SetColor(gui.NewQColor2(core.Qt__red))
			break
		case c.Warning:
			status = "Warning"
			b.SetColor(gui.NewQColor2(core.Qt__yellow))
			break
		case c.Depreated:
			status = "Deperated"
			b.SetColor(gui.NewQColor2(core.Qt__magenta))
			break
		case c.Unknown:
			status = "Unknown"
			b.SetColor(gui.NewQColor2(core.Qt__gray))
			break
		}
		item := widgets.NewQTreeWidgetItem2([]string{v, model, status}, 0)
		item.SetForeground(2, b)
		p.tw.InsertTopLevelItem(k, item)
	}

	p.tw.ExpandAll()
}
