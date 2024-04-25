package client

import (
	"context"
	"log"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Add to Server user's site credentials: login and password.
// If site created success on the server, it return new UUID of created site object.
func SiteAdd(conf config.Config, user oapi.User, siteURL, slogin, spw string) (site *oapi.Site, status int, err error) {

	// custom HTTP client
	hc := http.Client{}
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}
	//SetHeader("Authorization", config.AuthPrefix+*user.Jwt).
	//jwtf := NewExampleAuthProvider(config.AuthPrefix + *user.Jwt)

	// Create OAPI site object.
	site = &oapi.Site{Site: siteURL, Slogin: slogin, Spw: spw}
	resp, err := c.AddSiteGen(context.TODO(), *site, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+*user.Jwt)
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

	zap.S().Infoln(user.Jwt, user.Login, user.Password)
	return site, resp.StatusCode, nil
}

func addJWT(ctx context.Context, req *http.Request) error {

	return nil
}
