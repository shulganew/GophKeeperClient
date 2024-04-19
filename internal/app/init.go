package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/loginpw"
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

func InitModel(conf config.Config) tea.Model {
	// Load User from backup
	cSate := tui.MainMenu
	user, err := backup.LoadUser()
	if err != nil {
		zap.S().Infoln("Saved user not found.")
		cSate = tui.NotLoginMenu
	}
	zap.S().Infoln("Start menu: ", cSate)
	//
	// Menu: Init Not Login, state 0.
	//
	nl := states.NewNotLogin()
	// Login form, state 1.
	lf := states.NewLoginForm()
	// Register form - state 2.
	rf := states.NewRegisterForm()
	//
	// Menu: Main menu for loged in users. State 3.
	//
	mm := states.NewMainMenu()
	//
	// Menu: Save site's login and passwords. 4
	//
	lm := loginpw.NewLoginMenu()
	// List site's login and passwords 5
	ll := loginpw.NewListLogin()
	// Add site's login and passwords 6
	al := loginpw.NewAddLogin()

	return tui.Model{Conf: conf, CurrentState: cSate, States: []tui.State{&nl, &lf, &rf, &mm, &lm, &ll, &al}, User: user}
}
