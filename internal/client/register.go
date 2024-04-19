package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/entities"
	"go.uber.org/zap"
)

const authPrefix = "Bearer "

func UserReg(conf config.Config, login, email, pw string) (user *entities.User, status int, err error) {
	// Create user.
	user = &entities.User{Login: login, Email: email, Password: pw}
	// Encode to json.
	bodyUser := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(bodyUser).Encode(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Make requset to server.

	url := url.URL{Scheme: config.Shema, Host: conf.Address, Path: config.RegisterPath}
	zap.S().Debugln("Register URL: ", url)

	client := resty.New()
	res, err := client.R().
		SetBody(bodyUser).
		SetHeader("Content-Type", "application/json").
		Post(url.String())

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range res.Header() {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	zap.S().Debugln("Body: ", string(res.Body()))
	zap.S().Debugf("Status Code: %d\r\n", res.StatusCode)
	authHeader := res.Header().Get("Authorization")
	zap.S().Debugln("authHeader: ", authHeader)
	if strings.HasPrefix(authHeader, authPrefix) {
		jwt := authHeader[len(authPrefix):]
		user.JWT = sql.NullString{String: jwt, Valid: true}
	} else {
		user.JWT = sql.NullString{String: "", Valid: false}
	}

	zap.S().Infoln(user.JWT, user.Login, user.Password)
	return user, res.StatusCode(), nil
}
