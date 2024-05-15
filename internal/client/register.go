package client

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

const authPrefix = "Bearer "

func UserReg(c *oapi.Client, ctx context.Context, conf config.Config, login, pw, email string) (user *oapi.NewUser, jwt string, status int, secret string, err error) {
	// Create OAPI user.
	user = &oapi.NewUser{Login: login, Password: pw, Email: email}
	resp, err := c.CreateUser(ctx, *user)
	if err != nil {
		return nil, "", http.StatusInternalServerError, "", err
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

	// Read otp secret from body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, "", err
	}

	return user, jwt, resp.StatusCode, string(data), nil
}
