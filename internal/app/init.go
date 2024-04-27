package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/ccard"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/site"
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
		user = &backup.BackupData{}
	}

	zap.S().Debugln("Start menu: ", cSate)
	//
	// Menu: Init Not Login, state 0.
	//
	nl0 := states.NewNotLogin()
	// Login form, state 1.
	lf1 := states.NewLoginForm()
	// Register form - state 2.
	rf2 := states.NewRegisterForm()
	//
	// Menu: Main menu for loged in users. State 3.
	//
	mm3 := states.NewMainMenu()
	//
	// Menu: Save site's login and passwords. 4
	//
	lm4 := site.NewLoginMenu()
	// List site's login and passwords 5
	sl5 := site.NewSiteList()
	// Add site's login and passwords 6
	al6 := site.NewAddLogin()
	// TODO 7-10  temp chops
	stub7 := states.NewMainMenu()
	stub8 := states.NewMainMenu()
	stub9 := states.NewMainMenu()
	stub10 := states.NewMainMenu()
	// Add site's login and passwords 11
	ca11 := ccard.NewCardAdd()

	// TODO make transfer object and Model constructor
	return tui.Model{Conf: conf, User: &user.NewUser, JWT: user.JWT, CurrentState: cSate, States: []tui.State{&nl0, &lf1, &rf2, &mm3, &lm4, &sl5, &al6, &stub7, &stub8, &stub9, &stub10, &ca11}}
}
