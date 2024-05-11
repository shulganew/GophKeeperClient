package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/app/logging"

	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)

func main() {
	logging.InitLog()
	zap.S().Infoln("Start app")

	conf := config.InitConfig()
	initialModel := app.InitModel(*conf, buildVersion, buildDate)

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		zap.S().Errorln("could not start program:", err)
	}
}
