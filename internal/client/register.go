package client

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

const authPrefix = "Bearer "

func UserReg(ctx context.Context, conf config.Config, login, email, pw string) (user *oapi.User, status int, err error) {
	// custom HTTP client
	hc := http.Client{}

	// with a raw http.Response

	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}
	// Create OAPI user.
	user = &oapi.User{Login: login, Email: email, Password: pw}
	resp, err := c.UserRegGen(ctx, *user)
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
	authHeader := resp.Header.Get("Authorization")
	zap.S().Debugln("authHeader: ", authHeader)
	if strings.HasPrefix(authHeader, authPrefix) {
		jwtStr := authHeader[len(authPrefix):]
		user.Jwt = &jwtStr
	}
	zap.S().Infoln(user.Jwt, user.Login, user.Password)

	return user, resp.StatusCode, nil
}
