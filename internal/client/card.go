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

// Add to Server user's card credentials: login and password.
// If card created success on the server, it return new UUID of created card object.
func CardAdd(conf config.Config, jwt, def, ccn, cvv, exp, hld string) (ncard *oapi.NewCard, status int, err error) {

	// custom HTTP client
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(GetTLSClietn()))
	if err != nil {
		log.Fatal(err)
	}

	// Create OAPI card object.
	ncard = &oapi.NewCard{Definition: def, Ccn: ccn, Cvv: cvv, Exp: exp, Hld: hld}
	// Add saved jwt token for auth.
	resp, err := c.AddCard(context.TODO(), *ncard, func(ctx context.Context, req *http.Request) error {
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

	return ncard, resp.StatusCode, nil
}

// Retrive all cards credentials from the server.
func CardsList(conf config.Config, jwt string) (cards []oapi.Card, status int, err error) {

	// custom HTTP client
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(GetTLSClietn()))
	if err != nil {
		log.Fatal(err)
	}

	// Create OAPI card object.
	resp, err := c.ListCards(context.TODO(), func(ctx context.Context, req *http.Request) error {
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

	// Decode cards from body.
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&cards)
		if err != nil {
			zap.S().Errorln("Can't write to response in ListCard handler", err)
		}
	}

	return cards, resp.StatusCode, nil
}
