package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/entities"
	"go.uber.org/zap"
)

const RegisterPath = "/api/user/register"
const authPrefix = "Bearer "

func UserReg(conf config.Config, login, email, pw string) (user *entities.User, status int, err error) {
	// Create user.
	user = &entities.User{Login: login, Email: email, Password: pw}
	// Encode to json.
	reqBodyDel := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(reqBodyDel).Encode(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Make requset to server.
	client := &http.Client{}
	url := url.URL{Scheme: config.Shema, Host: conf.Address, Path: RegisterPath}
	zap.S().Debugln("Register URL: ", url)
	request, err := http.NewRequest(http.MethodPost, url.String(), reqBodyDel)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	request.Header.Add("Content-Type", "application/json")
	// Get response, check error.
	res, err := client.Do(request)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range res.Header {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	zap.S().Debugln("Body: ", string(body))
	zap.S().Debugf("Status Code: %d\r\n", res.StatusCode)
	authHeader := res.Header.Get("Authorization")
	zap.S().Debugln("authHeader: ", authHeader)
	if strings.HasPrefix(authHeader, authPrefix) {
		jwt := authHeader[len(authPrefix):]
		user.JWT = sql.NullString{String: jwt, Valid: true}
	} else {
		user.JWT = sql.NullString{String: "", Valid: false}
	}

	zap.S().Infoln(user.JWT, user.Login, user.Password)
	return user, res.StatusCode, nil
}
