package color

import "github.com/charmbracelet/lipgloss"

const NONE_COLOR = ""

type ColorFrame struct {
	Text       lipgloss.Color
	Background lipgloss.Color
}

type ColorScheme struct {
	Fkey            *ColorFrame
	BorderBase      *ColorFrame
	BorderHighlight *ColorFrame
	AppName         *ColorFrame
	TextBase        *ColorFrame
	TextBaseDark    *ColorFrame
	TextHighlight   *ColorFrame
	Title           *ColorFrame
	ServerStatus    map[string]*ColorFrame

	//input
	Placeholder   *ColorFrame
	ButtonFocused *ColorFrame
	ButtonBlurred *ColorFrame
	Help          *ColorFrame
	CursorHelp    *ColorFrame
}

func PinkAndrewScheme() *ColorScheme {
	return &ColorScheme{
		Fkey:            &ColorFrame{Text: lipgloss.Color("205")},
		BorderHighlight: &ColorFrame{Text: lipgloss.Color("200")},
		AppName:         &ColorFrame{Text: lipgloss.Color("205")},
		TextBase:        &ColorFrame{Text: lipgloss.Color("250")},
		TextBaseDark:    &ColorFrame{Text: lipgloss.Color("240")},
		TextHighlight:   &ColorFrame{Background: lipgloss.Color("200")},
		Title:           &ColorFrame{Text: lipgloss.Color("205"), Background: lipgloss.Color("240")},
		Placeholder:     &ColorFrame{Text: lipgloss.Color("245")},
		ButtonFocused:   &ColorFrame{Text: lipgloss.Color("205")},
		ButtonBlurred:   &ColorFrame{Text: lipgloss.Color("74")},
		Help:            &ColorFrame{Text: lipgloss.Color("74")},
		CursorHelp:      &ColorFrame{Text: lipgloss.Color("244")},
		ServerStatus: map[string]*ColorFrame{
			"connected":    {Text: lipgloss.Color("42")},
			"connecting":   {Text: lipgloss.Color("214")},
			"disconnected": {Text: lipgloss.Color("196")},
		},
	}
}

func OceanScheme() *ColorScheme {
	return &ColorScheme{
		Fkey:            &ColorFrame{Text: lipgloss.Color("33")},
		BorderBase:      &ColorFrame{Text: lipgloss.Color("238")},
		BorderHighlight: &ColorFrame{Text: lipgloss.Color("39")},
		AppName:         &ColorFrame{Text: lipgloss.Color("39")},
		TextBase:        &ColorFrame{Text: lipgloss.Color("245")},
		TextBaseDark:    &ColorFrame{Text: lipgloss.Color("238")},
		TextHighlight:   &ColorFrame{Background: lipgloss.Color("33")},
		Title:           &ColorFrame{Text: lipgloss.Color("33"), Background: lipgloss.Color("240")},
		Placeholder:     &ColorFrame{Text: lipgloss.Color("245")},
		ButtonFocused:   &ColorFrame{Text: lipgloss.Color("74")},
		ButtonBlurred:   &ColorFrame{Text: lipgloss.Color("238")},
		Help:            &ColorFrame{Text: lipgloss.Color("74")},
		CursorHelp:      &ColorFrame{Text: lipgloss.Color("244")},
		ServerStatus: map[string]*ColorFrame{
			"connected":    {Text: lipgloss.Color("42")},
			"connecting":   {Text: lipgloss.Color("226")},
			"disconnected": {Text: lipgloss.Color("196")},
		},
	}
}

var GColorScheme = OceanScheme()
