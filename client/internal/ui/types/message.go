package types

import (
	"andrew_chat/client/internal/domain"

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

// DeleteWindowMsg is a message sent to the window manager
// instructing it to delete the specified window.
type DeleteWindowMsg struct {
	Model tea.Model
}

// TerminateWindow is a signal indicating that the window
// should terminate its execution and stop running.
type TerminateWindow struct{
	Model tea.Model
}

type FilterMsg struct {
	Filter string
}

type ErrMsg struct {
	Text string
}

type ReloadMsg struct{}