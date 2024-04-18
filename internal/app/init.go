package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init zap logger as main logger.
func InitLog() zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"/tmp/gkc2.log"},
		ErrorOutputPaths: []string{"/tmp/gkc2.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	zapLogger := zap.Must(cfg.Build())
	zapLogger.Info("logger construction succeeded")
	zap.ReplaceGlobals(zapLogger)
	defer func() {
		_ = zapLogger.Sync()
	}()

	sugar := *zapLogger.Sugar()

	defer func() {
		_ = sugar.Sync()
	}()
	return sugar
}

func InitModel() tea.Model {

	// Init Not Login.
	nl := states.NotLogin{Choices: []string{"Log In", "Sign Up", "Get out"}}
	return tui.Model{CurrentState: 0, States: []tui.State{&nl}}
}
