package client

import (
	"context"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

func UserLogin(c *oapi.Client, login, pw, otp string) (user *oapi.NewUser, jwt string, status int, err error) {
	// Create OAPI user.
	user = &oapi.NewUser{Login: login, Password: pw, Otp: otp}
	resp, err := c.Login(context.Background(), *user)
	if err != nil {
		return nil, "", 0, err
	}

	// Print to log file for debug level.
	for k, v := range resp.Header {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	// Get JWT token and save to User
	authHeader := resp.Header.Get("Authorization")
	zap.S().Debugln("authHeader: ", authHeader)
	if strings.HasPrefix(authHeader, authPrefix) {
		jwt = authHeader[len(authPrefix):]
	}
	zap.S().Infoln(jwt, user.Login, user.Password)
	return user, jwt, resp.StatusCode, nil
}
