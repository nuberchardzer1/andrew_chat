package wm

import (
	// "andrew_chat/intenal/tui/windows"
	"andrew_chat/intenal/color"
	"andrew_chat/intenal/debug"
	"andrew_chat/intenal/ui/keys"
	"andrew_chat/intenal/ui/types"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const notFocused = types.PositionSentinel

// implements bubbletea.model

type WindowManager struct {
	windows map[types.Position]types.ComponentModel
	focus   types.Position
	width   int
	height  int

	stack            []types.ComponentModel //when add window to stack it grows when delete evict
	leftWindowWidth  int
	rightWindowWidth int
}

func (wm *WindowManager) getWindowWidthWithoutBorder(p types.Position) int {
	if wm.height <= 0 {
		return 0
	}
	if p.IsLeft() {
		return int(float32(wm.width)*0.2) - types.BorderSize
	} else if p.IsRight() {
		return int(float32(wm.width)*0.8) - types.BorderSize
	}
	panic("unknown position")
}

func (wm *WindowManager) getWindowHeightWithoutBorder(p types.Position) int {
	var count int

	if p.IsLeft() {
		if wm.getWindow(types.PositionTopLeft) != nil {
			count++
		}
		if wm.getWindow(types.PositionBotLeft) != nil {
			count++
		}
	} else if p.IsRight() {
		if wm.getWindow(types.PositionTopRight) != nil {
			count++
		}
		if wm.getWindow(types.PositionBotRight) != nil {
			count++
		}
	} else {
		panic(fmt.Sprintf("unexpected position %08b", p))
	}

	if count == 0 {
		return 0
	}

	totalAvailable := wm.height - types.BorderSize*(count-1)
	return totalAvailable / count
}

func (wm *WindowManager) removeWindow(win types.ComponentModel) {
	p := types.PositionSentinel
	for pos, w := range wm.windows {
		if w == win {
			p = pos
			break
		}
	}
	if p == types.PositionSentinel {
		panic("delete unknown model not allowed")
	}

	debug.DebugDump(debug.V, fmt.Sprintf("Remove window position: %d", p), win) 
	if wm.focus == p {
		wm.nextPosition()
	}
	delete(wm.windows, p)
	wm.updateWindows()
}

func (wm *WindowManager) setWindow(p types.Position, win types.ComponentModel) {
	if win == nil {
		panic("set nil window not allowed")
	}

	if p.IsBot() {
		if wm.getWindow(p.SetTop()) == nil {
			panic("bot hasnt top sibling")
		}
	}
	wm.windows[p] = win
	wm.focus = p

	debug.DebugDump(debug.V, fmt.Sprintf("SET WINDOW POSITION: %d", p), win)
	win.Init()
	wm.updateWindows()
}

func (wm *WindowManager) getWindow(p types.Position) tea.Model {
	if win, ok := wm.windows[p]; ok {
		return win
	}
	return nil
}

func (wm *WindowManager) nextPosition() {
	pos := wm.focus

	if pos.IsTop() {
		newPos := pos.SetBot()
		if wm.getWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsLeft() {
		newPos := pos.SetRight().SetTop()
		if wm.getWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsRight() {
		newPos := pos.SetLeft().SetTop()
		if wm.getWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsBot() {
		newPos := pos.SetTop()
		if wm.getWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	wm.focus = notFocused
}

func NewWM() *WindowManager {
	return &WindowManager{
		windows: make(map[types.Position]types.ComponentModel),
		focus:   notFocused,
	}
}

func (m *WindowManager) Init() tea.Cmd {
	return nil
}

func (wm *WindowManager) updateWindows() {
	for pos, win := range wm.windows {
		win.Update(
			tea.WindowSizeMsg{
				Height: wm.getWindowHeightWithoutBorder(pos),
				Width:  wm.getWindowWidthWithoutBorder(pos),
			},
		)
	}
}

func (m *WindowManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	debug.DebugDump(debug.VV, "UPDATE", m, msg)
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch {
		// These keys should exit the program.
		case key.Matches(msg, keys.Keys.Next):
			m.nextPosition()
			return m, nil
		default:
			if m.focus != notFocused {
				c := m.windows[m.focus]
				_, cmd := c.Update(msg)
				return m, cmd
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.leftWindowWidth = int(float32(m.width) * 0.2)
		m.rightWindowWidth = m.width - m.leftWindowWidth

		m.updateWindows()

	case types.CreateWindowMsg:
		m.setWindow(msg.Pos, msg.Model)

	case types.DeleteWindowMsg:
		debug.DebugDump(debug.V, "WM DeleteWindowMsg", msg)
		m.removeWindow(msg.Model)

	default:
		if m.focus != notFocused {
			c := m.windows[m.focus]
			_, cmd := c.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (wm *WindowManager) renderWindow(pos types.Position) string {
	win, ok := wm.windows[pos]
	if !ok {
		return ""
	}

	width := wm.getWindowWidthWithoutBorder(pos)
	height := wm.getWindowHeightWithoutBorder(pos)

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Top) // FIX

	if wm.focus == pos {
		style = style.BorderForeground(color.GColorScheme.BorderHighlight.Text)
	}

	return style.Render(win.View())
}

func (wm *WindowManager) View() string {
	if wm.focus == notFocused {
		return ""
	}

	wm.updateWindows() //hard code. i have no clue where is bug
	res := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Top,
			removeEmptyStrings(
				wm.renderWindow(types.PositionTopLeft),
				wm.renderWindow(types.PositionBotLeft),
			)...,
		),
		lipgloss.JoinVertical(
			lipgloss.Top,
			removeEmptyStrings(
				wm.renderWindow(types.PositionTopRight),
				wm.renderWindow(types.PositionBotRight),
			)...,
		),
	)
	debug.DebugDump(debug.VVV, "VIEW", wm, res)
	return res
}

func removeEmptyStrings(strs ...string) []string {
	n := 0
	for _, s := range strs {
		if s != "" {
			strs[n] = s
			n++
		}
	}
	return strs[:n]
}
