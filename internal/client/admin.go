package client

import (
	"context"
	"net/http"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
)

// Add to Server user's card credentials: login and password.
// If card created success on the server, it return new UUID of created card object.
func CreateNewEKey(c *oapi.Client, jwt string) (status int, err error) {
	resp, err := c.EKeyNew(context.TODO(), func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func CrateNewMaster(c *oapi.Client, jwt string, okey, nkey string) (status int, err error) {
	keys := oapi.Key{New: nkey, Old: okey}
	resp, err := c.NewMaster(context.TODO(), keys, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
