package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/entities"
	"go.uber.org/zap"
)

// Add to Server user's site credentials: login and password.
func SiteAdd(conf config.Config, user entities.User, siteURL, slogin, spw string) (status int, err error) {
	// Create JSON requset
	site := &entities.Site{SiteURL: siteURL, SLogin: slogin, SPw: spw}
	bodySite := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(bodySite).Encode(&site)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Make requset to server.
	url := url.URL{Scheme: config.Shema, Host: conf.Address, Path: config.SiteAddPath}
	zap.S().Debugln("Login URL: ", url)

	client := resty.New()
	res, err := client.R().
		SetBody(bodySite).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", config.AuthPrefix+user.JWT.String).
		Post(url.String())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range res.Header() {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	zap.S().Debugln("Body: ", string(res.Body()))
	zap.S().Debugf("Status Code: %d\r\n", res.StatusCode)

	return res.StatusCode(), nil

}
