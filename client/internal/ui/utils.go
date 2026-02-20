package ui

import (
	"andrew_chat/client/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

func NewDeleteCmd(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return types.DeleteWindowMsg{
			Model: model,
		}
	}
}

func NewCreateCmd(pos types.Position, model tea.Model, focus bool) tea.Cmd {
	return func() tea.Msg {
		return types.CreateWindowMsg{
			Pos:   pos,
			Model: model,
			Focus: focus,
		}
	}
}

func NewErrCmd(text string) tea.Cmd {
	return func() tea.Msg {
		return types.ErrMsg{
			Text: text,
		}
	}
}
