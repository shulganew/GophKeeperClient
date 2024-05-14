package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Add to Server user's card credentials: login and password.
// If card created success on the server, it return new UUID of created card object.
func CardAdd(c *oapi.Client, conf config.Config, jwt, def, ccn, cvv, exp, hld string) (ncard *oapi.NewCard, status int, err error) {

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

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	// Get JWT token and save to User
	return ncard, resp.StatusCode, nil
}

// Retrive all cards credentials from the server.
func CardsList(c *oapi.Client, conf config.Config, jwt string) (cards map[string]oapi.Card, status int, err error) {

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

// Update card by id.
func CardsUpdate(c *oapi.Client, conf config.Config, jwt string, cardID, def, ccn, cvv, exp, hld string) (status int, err error) {
	// Create OAPI card object.
	card := &oapi.Card{CardID: cardID, Definition: def, Ccn: ccn, Cvv: cvv, Exp: exp, Hld: hld}
	resp, err := c.UpdateCard(context.TODO(), *card, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	return resp.StatusCode, nil
}
