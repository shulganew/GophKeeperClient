package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/card"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/gfile"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/gtext"
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
	lm4 := site.NewSietMenu()
	// List site's login and passwords 5
	sl5 := site.NewSiteList()
	// Add site's login and passwords 6
	al6 := site.NewSiteAdd()

	siteU7 := states.NewMainMenu() // Site edit reserved
	// Card menu
	cm8 := card.NewCardMenu()
	ca9 := card.NewCardAdd()
	cl10 := card.NewCardList()

	// Text menu
	mg11 := gtext.NewGtextMenu()
	gta12 := gtext.NewGtextAdd()
	gtl13 := gtext.NewGtextList()

	// Text menu
	gm14 := gfile.NewGfileMenu()
	gm15 := gfile.NewFileAdd()
	gtl16 := gfile.NewGfileList()

	// TODO make transfer object and Model constructor
	states := []tui.State{nl0, lf1, rf2, mm3, lm4, sl5, al6, siteU7, cm8, ca9, cl10, mg11, gta12, gtl13, gm14, gm15, gtl16}
	return tui.Model{Conf: conf, User: &user.NewUser, JWT: user.JWT, CurrentState: cSate, States: states}
}
