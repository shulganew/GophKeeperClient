package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app"
	"github.com/shulganew/GophKeeperClient/internal/tui/tuimodel"
	"go.uber.org/zap"
)

func main() {
	app.InitLog()
	zap.S().Infoln("Start app")
	initialModel := tuimodel.Model{Choices: []string{"Log In", "Sign Up"}, Question: "You are not authorized, Log In or Sign Up:"}
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
