package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// implements bubbletea.model
type Dummy struct {
	width  int
	height int
}


func NewDummy() *Dummy {
	return &Dummy{}
}

func (m *Dummy) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *Dummy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *Dummy) View() string {
	content := lipgloss.NewStyle().
		Width(m.width).
		// accommodate header and footer
		Height(m.height).
		// Border(lipgloss.DoubleBorder()).
		Render()

	// Send the UI for rendering
	return content
}
