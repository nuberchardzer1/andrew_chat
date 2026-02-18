package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// implements bubbletea.model
type TextView struct {
	viewport viewport.Model
}

func NewTextView() *TextView {
	vp := viewport.New(0, 0)
	return &TextView{
		viewport: vp,
	}
}

func (m *TextView) Init() tea.Cmd {
	return nil
}

func (m *TextView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *TextView) View() string {
	// Send the UI for rendering
	return m.viewport.View()
}

func (m *TextView) SetContent(content string) {
	m.viewport.SetContent(content)
}
