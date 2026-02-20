package ui

import (
	"andrew_chat/client/internal/ui/keys"
	"andrew_chat/client/internal/ui/types"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Option struct {
	Name   string
	Action func() tea.Cmd
}

func (o Option) Title() string       { return o.Name }
func (o Option) Description() string { return "" }
func (o Option) FilterValue() string { return o.Name }

type ControlPane struct {
	width    int
	height   int
	listview list.Model
}

func NewControlPane(opts []Option) *ControlPane {
	var items []list.Item
	for _, opt := range opts {
		items = append(items, opt)
	}

	delegate := list.NewDefaultDelegate()
	listview := list.New(items, delegate, 0, 0)

	listview.FilterInput.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")).Bold(true)

	listview.SetShowTitle(false)
	listview.SetShowStatusBar(false)
	listview.SetShowHelp(false)
	listview.SetShowPagination(false)
	listview.SetFilteringEnabled(false)

	return &ControlPane{
		listview: listview,
	}
}

func (m *ControlPane) Init() tea.Cmd {
	return nil
}

func (m *ControlPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Choose):
			selectedItem := m.listview.SelectedItem()
			if selectedItem == nil {
				return m, nil
			}
			opt := selectedItem.(Option)

			return m, tea.Sequence(opt.Action(), NewDeleteCmd(m))
		case key.Matches(msg, keys.Keys.Close):
			return m, NewDeleteCmd(m)
		}
	case tea.WindowSizeMsg:
		m.listview.SetSize(msg.Width, msg.Height)
		return m, nil
	case types.FilterMsg:
		m.listview.SetFilterText(msg.Filter)
		return m, nil
	}

	m.listview, cmd = m.listview.Update(msg)
	return m, cmd
}

func (m *ControlPane) View() string {
	return m.listview.View()
}
