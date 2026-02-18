package ui

import (
	"andrew_chat/intenal/color"
	"andrew_chat/intenal/ui/types"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// wrapper for list.Model
type List struct {
	listview list.Model
}

func (m *List) Init() tea.Cmd {
	return nil
}

func (m *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.listview.SetSize(msg.Width, msg.Height)
		return m, nil

	case types.FilterMsg:
		m.listview.SetFilterText(msg.Filter)
		return m, nil

	default:
		var newList list.Model
		newList, cmd = m.listview.Update(msg)
		m.listview = newList
		return m, cmd
	}
}

func (m *List) View() string {
	return m.listview.View()
}

func NewList(items []list.Item, delegate list.ItemDelegate,
	width int, height int) *List {
	listview := list.New(items, delegate, width, height)

	// change Filter and cursor foregorund colors
	listview.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Bold(true)
	listview.FilterInput.Cursor.Style = lipgloss.NewStyle().
		Background(color.GColorScheme.TextHighlight.Background).
		Foreground(color.GColorScheme.Title.Text)

	// listview.Title = "Servers"
	listview.SetShowTitle(false)
	listview.SetStatusBarItemName("server", "servers")
	listview.SetShowStatusBar(false)
	listview.SetShowHelp(false)
	listview.SetShowPagination(false)
	listview.SetFilteringEnabled(false)
	// listview.SetShowFilter(true)

	var list List
	list.listview = listview
	return &list
}

func (m *List) SelectedItem() list.Item {
	return m.listview.SelectedItem()
}
