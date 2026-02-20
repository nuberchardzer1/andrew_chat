package keys

import "github.com/charmbracelet/bubbles/key"

type AppKeys struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Next   key.Binding
	Choose key.Binding
	Close  key.Binding
}

var Keys = AppKeys{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c", "q"),
	),
	Choose: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter", " "),
	),
	Close: key.NewBinding(
		key.WithKeys("esc", "esc"),
		key.WithHelp("enter", "esc"),
	),
}
