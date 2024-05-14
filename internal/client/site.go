package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Add to Server user's site credentials: login and password.
// If site created success on the server, it return new UUID of created site object.
func SiteAdd(c *oapi.Client, jwt, def, siteURL, slogin, spw string) (nsite *oapi.NewSite, status int, err error) {
	// Create OAPI site object.
	nsite = &oapi.NewSite{Definition: def, Site: siteURL, Slogin: slogin, Spw: spw}
	resp, err := c.AddSite(context.TODO(), *nsite, func(ctx context.Context, req *http.Request) error {
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	zap.S().Debugln("Body: ", string(body))
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	return nsite, resp.StatusCode, nil
}

// Retrive all sites credentials from the server.
func SiteList(c *oapi.Client, jwt string) (sites map[string]oapi.Site, status int, err error) {
	// Create OAPI site object.
	resp, err := c.ListSites(context.TODO(), func(ctx context.Context, req *http.Request) error {
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

// Site update by id.
func SiteUpdate(c *oapi.Client, jwt string, siteID, def, siteURL, slogin, spw string) (status int, err error) {
	// Create OAPI site object.
	site := &oapi.Site{SiteID: siteID, Definition: def, Site: siteURL, Slogin: slogin, Spw: spw}
	resp, err := c.UpdateSite(context.TODO(), *site, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)
	return resp.StatusCode, nil
}
