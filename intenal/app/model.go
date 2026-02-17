package app

import (
	"andrew_chat/intenal/color"
	"andrew_chat/intenal/config"
	"andrew_chat/intenal/server"
	"andrew_chat/intenal/ui"
	uisrv "andrew_chat/intenal/ui/server"
	"andrew_chat/intenal/ui/types"
	wm "andrew_chat/intenal/ui/window_manager"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	xPadding = 1
	yPadding = 0
)

// =============================================================================
// Util functions
// =============================================================================

func matchServerStatus(s int) string {
	var text string
	switch s {
	case server.StatusConnected:
		text = "connected"
	case server.StatusConnecting:
		text = "connecting"
	case server.StatusDisconnected:
		text = "disconnected"
	default:
		panic("unknown server status")
	}
	return text
}

// =============================================================================
// MainModel
// =============================================================================

type MainModel struct {
	wm    *wm.WindowManager
	focus int

	width  int
	height int

	server string
	status string
}

func NewMainModel() *MainModel {
	wm := wm.NewWM()
	return &MainModel{
		wm:           wm,
		status:       "disconnected",
	
	}
}


func navCmd(pos types.Position, model types.ComponentModel) tea.Cmd {
	return func() tea.Msg {
		return types.CreateWindowMsg{Pos: pos, Model: model}
	}
}

func (m *MainModel) Init() tea.Cmd {
	return tea.Sequence(
		// navCmd(types.PositionBotLeft, components.NewDummy()),
		navCmd(types.PositionTopRight, ui.NewDummy()),		
		navCmd(types.PositionTopLeft, uisrv.NewServer(config.GetServers())),
		// navCmd(types.PositionBotRight, components.NewDummy()),
	)
}

var testDummy = ui.NewDummy()

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF1:
		case tea.KeyF2:
			m.wm.Update(
				types.CreateWindowMsg{
					Pos:   types.PositionTopLeft,
					Model: uisrv.NewServer(config.GetServers()),
				},
				)
		case tea.KeyF3:
			m.wm.Update(
				types.CreateWindowMsg{
					Pos:   types.PositionTopRight,
					Model: testDummy,
				},
			)
		case tea.KeyF5:
			m.wm.Update(
				types.DeleteWindowMsg{
					Model: testDummy,
				},
			)
		case tea.KeyF10, tea.KeyCtrlC:
			return m, tea.Quit
		default:
			_, cmd := m.wm.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.wm.Update(tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height - types.HeaderHeight - types.FooterHeight,
		})
	case types.ServerMsg:
		if msg.Success {
			m.status = matchServerStatus(msg.Status)
			m.server = msg.Server.Address
		}
	default:
		_, cmd := m.wm.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *MainModel) renderStatus() string {
	color := color.GColorScheme.ServerStatus[m.status]
	var text string

	switch m.status {
	case "connected":
		text = "Connected"
	case "connecting":
		text = "Connecting"
	case "disconnected":
		text = "Disconnected"
	}

	return lipgloss.NewStyle().
		Foreground(color.Text).
		Render("● " + text)
}

func (m *MainModel) renderHeader() string {
	var leftSide string
	leftSide += lipgloss.NewStyle().
		Bold(true).
		Foreground(color.GColorScheme.AppName.Text).
		Render("AndrewChat" + "  ")

	if m.status == "connected" {
		leftSide += lipgloss.NewStyle().
			Foreground(color.GColorScheme.TextBase.Text).
			Render("Server: " + m.server)
	}

	status := m.renderStatus()

	space := m.width - lipgloss.Width(leftSide) - lipgloss.Width(status) - 1
	if space < 0 {
		space = 0
	}
	spacesBetween := lipgloss.NewStyle().
		Render(strings.Repeat(" ", space))
	row := leftSide + spacesBetween + status

	return lipgloss.NewStyle().
		Render(row)
}

func (m *MainModel) renderFooter() string {
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color.GColorScheme.Fkey.Text))

	itemStyle := lipgloss.NewStyle().
		Padding(0, 1)

	separator := lipgloss.NewStyle().
		Foreground(color.GColorScheme.TextBaseDark.Text).
		Render("│")

	// The footer
	footerItems := []string{
		itemStyle.Render(keyStyle.Render("F1") + " Setup  "),
		itemStyle.Render(keyStyle.Render("F2") + " Servers"),
		itemStyle.Render(keyStyle.Render("F3") + " Chats  "),
		itemStyle.Render(keyStyle.Render("F5") + " Refresh"),
		itemStyle.Render(keyStyle.Render("F9") + " Settings"),
		itemStyle.Render(keyStyle.Render("F10") + " Exit  "),
	}

	return lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Center,
			footerItems[0],
			separator,
			footerItems[1],
			separator,
			footerItems[2],
			separator,
			footerItems[3],
			separator,
			footerItems[4],
			separator,
			footerItems[5],
		))

}

func (m *MainModel) View() string {
	if m.height == 0 || m.width == 0 {
		return ""
	}
	// The header
	header := m.renderHeader()

	// Iterate over our choices
	content := m.wm.View()

	footer := m.renderFooter()
	// Send the UI for rendering
	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
