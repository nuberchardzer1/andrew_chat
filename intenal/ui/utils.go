package ui

import (
	"andrew_chat/intenal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

func NewDeleteCmd(model types.ComponentModel)tea.Cmd{
	return func() tea.Msg {
		return types.DeleteWindowMsg{
					Model: model,
				}
	}
}

func NewCreateCmd(pos types.Position, model types.ComponentModel)tea.Cmd{
	return func() tea.Msg {
		return types.CreateWindowMsg{
					Pos: pos,
					Model: model,
				}
	}
}

func NewErrCmd(text string)tea.Cmd{
	return func() tea.Msg {
		return types.ErrMsg{
					Text: text,
				}
	}
}