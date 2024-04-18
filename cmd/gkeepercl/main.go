package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app"
	"github.com/shulganew/GophKeeperClient/internal/app/config"

	"go.uber.org/zap"
)

func main() {
	app.InitLog()
	zap.S().Infoln("Start app")

	conf := config.InitConfig()
	initialModel := app.InitModel(*conf)

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		zap.S().Errorln("could not start program:", err)
	}
}
