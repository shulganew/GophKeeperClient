package client

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Add to Server user's gtext credentials: login and password.
// If gtext created success on the server, it return new UUID of created gtext object.
func GtextAdd(c *oapi.Client, conf config.Config, jwt, text string) (ngtext *oapi.NewGtext, status int, err error) {
	// Create OAPI gtext object.
	ngtext = &oapi.NewGtext{Definition: getDef(&text), Note: text}
	// Add saved jwt token for auth.
	resp, err := c.AddGtext(context.TODO(), *ngtext, func(ctx context.Context, req *http.Request) error {
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

	return ngtext, resp.StatusCode, nil
}

// Retrive all gtexts credentials from the server.
func GtextList(c *oapi.Client, conf config.Config, jwt string) (gtexts map[string]oapi.Gtext, status int, err error) {

	// Create OAPI gtext object.
	resp, err := c.ListGtexts(context.TODO(), func(ctx context.Context, req *http.Request) error {
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

	// Decode gtexts from body.
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&gtexts)
		if err != nil {
			zap.S().Errorln("Can't write to response in Listgtext handler", err)
		}
	}

	return gtexts, resp.StatusCode, nil
}

// Site update by id.
func GtextUpdate(c *oapi.Client, conf config.Config, jwt string, gtextID string, text string) (status int, err error) {
	// Create OAPI text object.
	gtext := &oapi.Gtext{GtextID: gtextID, Definition: getDef(&text), Note: text}
	resp, err := c.UpdateGtext(context.TODO(), *gtext, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)
	return resp.StatusCode, nil
}

// Return first sentence of the text.
func getDef(text *string) string {
	scanner := bufio.NewScanner(strings.NewReader(*text))
	if scanner.Scan() {
		return scanner.Text()
	}

	return "No header note."
}
