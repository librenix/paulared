package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var color = termenv.ColorProfile().Color

type errMsg error

type model struct {
	state    int
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
	if m.err != nil {
		return m.err.Error()
	}
	s := termenv.String(m.spinner.View()).Foreground(color("214")).String()
	title := fmt.Sprintf("\n\n   %s Checking hardware compatibility ...\n\n", s)
	page := ""
	//Compatibility()
	button := termenv.String("  Next  ").Background(color("214")).String()
	output := fmt.Sprintf("%s\n%s\n%s", title, page, button)
	if m.quitting {
		return output + "\n"
	}
	return output
}

func Start() *tea.Program {
	return tea.NewProgram(initialModel())
}
