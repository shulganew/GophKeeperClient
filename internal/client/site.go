package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Add to Server user's site credentials: login and password.
// If site created success on the server, it return new UUID of created site object.
func SiteAdd(conf config.Config, user oapi.NewUser, jwt, siteURL, slogin, spw string) (site *oapi.NewSite, status int, err error) {

	// custom HTTP client
	hc := http.Client{}
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}

	// Create OAPI site object.
	site = &oapi.NewSite{Site: siteURL, Slogin: slogin, Spw: spw}
	resp, err := c.AddSite(context.TODO(), *site, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range resp.Header {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	zap.S().Debugln("Body: ", resp.Body)
	zap.S().Debugln("Body: ", resp.Body)
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	// Get JWT token and save to User

	zap.S().Infoln(jwt, user.Login, user.Password)
	return site, resp.StatusCode, nil
}

// Retrive all sites credentials from the server.
func SiteList(conf config.Config, jwt string) (sites []oapi.Site, status int, err error) {

	// custom HTTP client
	hc := http.Client{}
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}

	// Create OAPI site object.
	resp, err := c.ListSite(context.TODO(), func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range resp.Header {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	zap.S().Debugln("Body: ", resp.Body)
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	// Get JWT token and save to User

	// Decode sites from body.
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&sites)
		if err != nil {
			zap.S().Errorln("Can't write to response in ListSite handler", err)
		}
	}

	return sites, resp.StatusCode, nil
}
