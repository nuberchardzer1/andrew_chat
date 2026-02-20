package ui

import (
	"andrew_chat/client/internal/color"
	"andrew_chat/client/internal/ui/keys"
	"andrew_chat/client/internal/ui/types"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(color.GColorScheme.ButtonFocused.Text)
	blurredStyle        = lipgloss.NewStyle().Foreground(color.GColorScheme.ButtonBlurred.Text)
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = lipgloss.NewStyle().Foreground(color.GColorScheme.Help.Text)
	PlaceholderStyle    = lipgloss.NewStyle().Foreground(color.GColorScheme.Placeholder.Text)
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(color.GColorScheme.CursorHelp.Text)

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type InputFormAction func([]types.InputFieldValue) tea.Cmd

type inputField struct {
	types.InputFieldSpec
	text textinput.Model
}
type InputFormModel struct {
	action     InputFormAction
	prompt     string
	focusIndex int
	cursorMode cursor.Mode
	fields     []inputField
}

func NewInputFormModel(prompt string, specs []types.InputFieldSpec,
	action InputFormAction) *InputFormModel {

	fields := make([]inputField, len(specs))
	for i, f := range specs {
		t := textinput.New()
		t.Placeholder = f.Placeholder
		t.Prompt = f.Title + ": "
		t.Focus()
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle
		t.Blur()
		t.PlaceholderStyle = PlaceholderStyle
		fields[i] = inputField{
			InputFieldSpec: f,
			text:           t,
		}
	}
	return &InputFormModel{
		action:     action,
		prompt:     prompt,
		focusIndex: 0,
		cursorMode: cursor.CursorBlink,
		fields:     fields,
	}
}

func (m *InputFormModel) Init() tea.Cmd {
	m.updateFocus()
	return textinput.Blink
}

func (m *InputFormModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.fields))
	for i := range m.fields {
		if i == m.focusIndex {
			cmds[i] = m.fields[i].text.Focus()
			m.fields[i].text.PromptStyle = focusedStyle
			m.fields[i].text.TextStyle = focusedStyle
		} else {
			m.fields[i].text.Blur()
			m.fields[i].text.PromptStyle = noStyle
			m.fields[i].text.TextStyle = noStyle
		}
	}
	return tea.Batch(cmds...)
}

func (m *InputFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Close):
			return m, NewDeleteCmd(m)
		case key.Matches(msg, keys.Keys.Up):
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			return m, m.updateFocus()
		case key.Matches(msg, keys.Keys.Down):
			if m.focusIndex < len(m.fields) {
				m.focusIndex++
			}
			return m, m.updateFocus()
		case key.Matches(msg, keys.Keys.Choose):
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.fields) {
				values := make([]types.InputFieldValue, len(m.fields))
				for i, f := range m.fields {
					values[i] = types.InputFieldValue{
						Name:  f.Name,
						Value: f.text.Value(),
					}
				}

				cmd := m.action(values)
				cmds = append(cmds, cmd, NewDeleteCmd(m))
				return m, tea.Batch(cmds...)
			}
			if s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex > len(m.fields) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.fields)
			}
			return m, m.updateFocus()
		}
	}
	return m, m.updateInputs(msg)
}

func (m *InputFormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.fields))
	for i := range m.fields {
		m.fields[i].text, cmds[i] = m.fields[i].text.Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *InputFormModel) View() string {
	var b strings.Builder
	for i := range m.fields {
		b.WriteString(m.fields[i].text.View())
		if i < len(m.fields)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.focusIndex == len(m.fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	return b.String()
}
