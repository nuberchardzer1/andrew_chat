package types

import (
	"andrew_chat/intenal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerMsg struct {
	Status int // "connect", "disconnect", "select"
	Server domain.Server
}

// The message is an instruction to the window manager where to place the window.
type CreateWindowMsg struct {
	Pos   Position
	Model tea.Model
	// focus on add or ignore
	Focus bool
}

// The message is an instruction to the window manager to delete the window.
type DeleteWindowMsg struct {
	Model tea.Model
}

type FilterMsg struct {
	Filter string
}

type ErrMsg struct {
	Text string
}
