package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shulganew/GophKeeperClient/internal/app/backup"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/shulganew/GophKeeperClient/internal/tui"
	"github.com/shulganew/GophKeeperClient/internal/tui/states"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/card"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/gfile"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/gtext"
	"github.com/shulganew/GophKeeperClient/internal/tui/states/site"
	"go.uber.org/zap"
)

func InitModel(conf config.Config, version, date string) tea.Model {
	// Client with TLS session.
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(GetTLSClietn(conf.SertPath)))
	if err != nil {
		log.Fatal(err)
	}

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
	nl0 := states.NewNotLogin(version, date)
	// Login form, state 1.
	lf1 := states.NewLoginForm()
	// Register form - state 2.
	rf2 := states.NewRegisterForm()
	//
	// Menu: Main menu for loged in users. State 3.
	//
	mm3 := states.NewMainMenu(version, date)
	//
	// Menu: Save site's login and passwords. 4
	//
	lm4 := site.NewSietMenu(version, date)
	// List site's login and passwords 5
	sl5 := site.NewSiteList()
	// Add site's login and passwords 6
	al6 := site.NewSiteAdd()

	siteU7 := site.NewSiteUpdate()
	// Card menu
	cm8 := card.NewCardMenu()
	ca9 := card.NewCardAdd()
	cl10 := card.NewCardList()
	cu11 := card.NewCardUpdate()

	// Text menu
	mg12 := gtext.NewGtextMenu(version, date)
	gta13 := gtext.NewGtextAdd()
	gtl14 := gtext.NewGtextList()
	gtup15 := gtext.NewGtextUpdate()

	// Text menu
	gm16 := gfile.NewGfileMenu(version, date)
	gm17 := gfile.NewFileAdd()
	gtl18 := gfile.NewGfileList(conf.FileSavingPath)

	states := []tui.State{nl0, lf1, rf2, mm3, lm4, sl5, al6, siteU7, cm8, ca9, cl10, cu11, mg12, gta13, gtl14, gtup15, gm16, gm17, gtl18}
	return tui.Model{Conf: conf, Client: c, User: &user.NewUser, JWT: user.JWT, CurrentState: cSate, States: states}
}

// Return http Client with TLS session config.
func GetTLSClietn(sertPath string) *http.Client {
	caCertf, _ := os.ReadFile(sertPath)
	rootCAs, _ := x509.SystemCertPool()
	// handle case where rootCAs == nil and create an empty pool...
	if ok := rootCAs.AppendCertsFromPEM(caCertf); !ok {
		zap.S().Infoln("Can't load trasted sertifecate!")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            rootCAs,
	}

	hc := &http.Client{
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}
	return hc
}
