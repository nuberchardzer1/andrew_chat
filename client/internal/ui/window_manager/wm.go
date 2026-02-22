package wm

import (
	"andrew_chat/client/internal/color"
	"andrew_chat/client/internal/debug"
	"andrew_chat/client/internal/ui"
	"andrew_chat/client/internal/ui/keys"
	"andrew_chat/client/internal/ui/types"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	leftPanelWidthRatio  = 0.2
	rightPanelWidthRatio = 0.8
	notFocused = types.PositionSentinel
)

type window struct {
	model tea.Model
	pos   types.Position
}

func newWindow(pos types.Position, m tea.Model) window {
	return window{
		m, pos,
	}
}

// implements bubbletea.model
type WindowManager struct {
	//visible windows
	windows map[types.Position]tea.Model
	focus   types.Position
	width   int
	height  int

	//TODO: struct for stack
	stack            []window
	leftWindowWidth  int
	rightWindowWidth int
}

func (wm *WindowManager) getWindowWidthWithoutBorder(p types.Position) int {
	if wm.height <= 0 {
		return 0
	}
	if p.IsLeft() {
		return int(float32(wm.width)*leftPanelWidthRatio) - types.BorderSize
	} else if p.IsRight() {
		return int(float32(wm.width)*rightPanelWidthRatio) - types.BorderSize
	}
	panic("unknown position")
}

func (wm *WindowManager) getWindowHeightWithoutBorder(p types.Position) int {
	var count int

	if p.IsLeft() {
		if wm.GetWindow(types.PositionTopLeft) != nil {
			count++
		}
		if wm.GetWindow(types.PositionBotLeft) != nil {
			count++
		}
	} else if p.IsRight() {
		if wm.GetWindow(types.PositionTopRight) != nil {
			count++
		}
		if wm.GetWindow(types.PositionBotRight) != nil {
			count++
		}
	} else {
		panic(fmt.Sprintf("unexpected position %08b", p))
	}

	if count == 0 {
		return wm.height - types.BorderSize
	}

	totalAvailable := wm.height - types.BorderSize*(count)
	return totalAvailable / count
}

func (wm *WindowManager) CloseWindow(win tea.Model)tea.Cmd{
	debug.DebugDump(debug.V, "Remove get", win)
	p := types.PositionSentinel
	for pos, w := range wm.windows {
		if w == win {
			p = pos
			break
		}
	}
	if p == types.PositionSentinel {
		// slog.Warn("try delete unknown window")
		return nil
	}

	debug.DebugDump(debug.V, fmt.Sprintf("Remove window pos: %d", p), win)
	delete(wm.windows, p)

	idx := -1

	for i := range wm.stack {
		if wm.stack[i].model == win {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic("inconsistent stack and windows")
	}

	wm.stack = append(wm.stack[:idx], wm.stack[idx+1:]...)
	if len(wm.stack) == 0 {
		wm.focus = notFocused
	} else if wm.focus == p {
		wm.focus = wm.stack[len(wm.stack)-1].pos
	}

	return ui.NewReloadCmd()
}

func (wm *WindowManager) AddWindow(p types.Position, win tea.Model, focus bool)tea.Cmd{
	if win == nil {
		panic("set nil window not allowed")
	}

	cmds := make([]tea.Cmd, 0)

	//window already occupied
	if oldWin, ok := wm.windows[p]; ok{
		debug.DebugDump(debug.V, "OLD WIN", oldWin, ok)
		cmds = append(cmds, ui.NewDeleteCmd(oldWin))
	}

	wm.windows[p] = win

	if focus {
		wm.stack = append(wm.stack, newWindow(p, win))
		wm.focus = p
	} else {
		newStack := make([]window, len(wm.stack)+1)
		newStack[0] = newWindow(p, win)
		copy(newStack[1:], wm.stack)
		wm.stack = newStack
		if wm.focus == notFocused {
			wm.focus = p
		}
	}

	cmds = append(cmds, ui.NewReloadCmd())
	debug.DebugDump(debug.V, fmt.Sprintf("Add window pos: %d, focus: %v", p, focus), win)
	return tea.Batch(cmds...)
}

func (wm *WindowManager) GetWindow(p types.Position) tea.Model {
	if win, ok := wm.windows[p]; ok {
		return win
	}
	return nil
}

func (wm *WindowManager) nextPosition() {
	pos := wm.focus
	defer func() {
		if wm.focus == notFocused {
			return
		}

		win := wm.windows[wm.focus]
		size := len(wm.stack)
		for i := size - 1; i >= 0; i-- {
			if wm.stack[i].model == win {
				wm.stack[i], wm.stack[size-1] = wm.stack[size-1], wm.stack[i]
				return
			}
		}
		panic("inconsistent stack and windows")
	}()

	if pos.IsTop() {
		newPos := pos.SetBot()
		if wm.GetWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsLeft() {
		newPos := pos.SetRight().SetTop()
		if wm.GetWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsRight() {
		newPos := pos.SetLeft().SetTop()
		if wm.GetWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}

	if pos.IsBot() {
		newPos := pos.SetTop()
		if wm.GetWindow(newPos) != nil {
			wm.focus = newPos
			return
		}
	}
}

func NewWM() *WindowManager {
	return &WindowManager{
		windows: make(map[types.Position]tea.Model),
		focus:   notFocused,
	}
}

func (m *WindowManager) Init() tea.Cmd {
	return nil
}

func (wm *WindowManager) updateWindows() {
	for pos, win := range wm.windows {
		updated, _ := win.Update(
			tea.WindowSizeMsg{
				Height: wm.getWindowHeightWithoutBorder(pos),
				Width:  wm.getWindowWidthWithoutBorder(pos),
			},
		)
		wm.windows[pos] = updated
	}
}

func (m *WindowManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	debug.DebugDump(debug.VV, "UPDATE", m, msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch {
		// These keys should exit the program.
		case key.Matches(msg, keys.Keys.Next):
			m.nextPosition()
			return m, nil

		case msg.String() == "ctrl+d":
			if m.focus == notFocused {
				return m, nil
			}
			return m, ui.NewDeleteCmd(m.stack[len(m.stack)-1].model)
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

		return m, ui.NewReloadCmd()

	case types.CreateWindowMsg:
		cmd := m.AddWindow(msg.Pos, msg.Model, msg.Focus)
		return m, tea.Sequence(msg.Model.Init(), cmd, ui.NewReloadCmd()) //need m.updateWindows()?
		
	case types.DeleteWindowMsg:
		debug.DebugDump(debug.V, "WM DeleteWindowMsg", msg)
		cmd := m.CloseWindow(msg.Model)
		cmds = append(cmds, cmd)
		_, cmd = msg.Model.Update(types.TerminateWindow{}) 
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case types.ReloadMsg:
		m.updateWindows()
		return m, nil

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
		Border(lipgloss.RoundedBorder()).
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
