package main

import (
	"andrew_chat/intenal/app"
	"andrew_chat/intenal/config"
	debug "andrew_chat/intenal/debug"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const logFile = "log.txt"

func main() {
	fmt.Println("andrew chat started")
	config.InitConfig("example.json")

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	debug.SetupLogger(file)
	// logger.AndrewLogDebugDump(file)
	app := app.NewApp()

	p := tea.NewProgram(app.MainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
