package client

import (
	"context"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Site update by id.
func Delete(c *oapi.Client, conf config.Config, jwt, secretID string) (status int, err error) {
	// Create OAPI site object.
	resp, err := c.DelAny(context.TODO(), secretID, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	return resp.StatusCode, nil
}
