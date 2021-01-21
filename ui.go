package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var color = termenv.ColorProfile().Color

type errMsg error

type model struct {
	state    uint
	choose   uint
	done     bool
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() model {
	s := spinner.NewModel()
	s.Spinner = spinner.Points
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.done {
				if m.state > 6 {
					return m, tea.Quit
				} else {
					m.state++
				}
			}
			return m, nil
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model) View() string {
	var s, banner, title, page, button string

	if m.err != nil {
		return m.err.Error()
	}
	s = termenv.String(m.spinner.View()).Foreground(color("214")).String()
	switch m.state {
	case 0:
		title = "Checking hardware compatibility ..."
		page, _ = Compatibility()
		break
	case 1:
		title = "Choose a dmg file"
		break
	case 2:
		title = "Choose a bootloader"
		break
	case 3:
		title = "Select target device"
		break
	case 4:
		title = "Downloading, please wait ..."
		break
	case 5:
		title = "Flashing, please don't eject your device ..."
		break
	case 6:
		title = "All done"
		m.quitting = true
		break
	}
	if m.done {
		banner = fmt.Sprintf("\n\n   ✓ %s \n\n", title)
	} else {
		banner = fmt.Sprintf("\n\n   %s %s \n\n", s, title)
	}
	page = ""
	if m.done {
		button = termenv.String("  Next  ").Background(color("214")).String()
	}
	output := fmt.Sprintf("%s\n%s\n%s", banner, page, button)
	if m.quitting {
		return output + "\n"
	}
	return output
}

func SetupUI() *tea.Program {
	return tea.NewProgram(initialModel())
}
