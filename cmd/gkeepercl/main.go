package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app"

	"go.uber.org/zap"
)

func main() {
	app.InitLog()
	zap.S().Infoln("Start app")

	initialModel := app.InitModel()

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		zap.S().Errorln("could not start program:", err)
	}
}
