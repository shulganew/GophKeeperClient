package main

import (
	"github.com/shulganew/GophKeeperClient/internal/app"
	"go.uber.org/zap"
)

func main() {
	app.InitLog()
	zap.S().Infoln("Start app")
}
