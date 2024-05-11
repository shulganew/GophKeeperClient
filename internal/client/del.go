package client

import (
	"context"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Sectet update by id (Site, Card, Text).
func DeleteAny(c *oapi.Client, jwt, secretID string) (status int, err error) {
	// Create OAPI site object.
	resp, err := c.DelAny(context.TODO(), secretID, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return resp.StatusCode, err
	}
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)
	return resp.StatusCode, nil
}

// Delete file from DB and Storage.
func DeleteFile(c *oapi.Client, jwt, fileID string) (status int, err error) {
	// Create OAPI site object.
	resp, err := c.DelGfile(context.TODO(), fileID, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return resp.StatusCode, err
	}
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)
	return resp.StatusCode, nil
}
