package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alecthomas/units"
	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

// Upload file metadata to server.
func FileAdd(c *oapi.Client, jwt, def, fPath string) (gfile *oapi.Gfile, status int, err error) {
	// Loading file form os
	file, err := os.Open(fPath)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Get file size.
	st, err := file.Stat()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// File size constrain 30 MIB.
	if st.Size() > int64(units.Mebibyte*30) {
		zap.S().Errorln("File too big, size less 30MiB.")
		return nil, 0, errors.New("file too big, size less 30MiB")
	}

	// Create nfile
	nfile := oapi.NewGfile{Definition: def, Fname: filepath.Base(file.Name()), Size: st.Size()}

	// Encode nfile as metadata to binary
	var md bytes.Buffer
	err = json.NewEncoder(&md).Encode(&nfile)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.AddGfileWithBody(context.TODO(), "application/json", &md, func(ctx context.Context, req *http.Request) error {
		// Add saved jwt token for auth.
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	status = resp.StatusCode
	// Decode gfile from body.
	if status == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&gfile)
		if err != nil {
			zap.S().Errorln("Can't write to response in Listgfile handler", err)
			return nil, http.StatusInternalServerError, err
		}

	}

	return
}

// Upload file  to server.
func FileUpload(c *oapi.Client, jwt, fPath, fileID string) (status int, err error) {
	// Loading file form os.
	file, err := os.Open(fPath)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	resp, err := c.UploadGfileWithBody(context.TODO(), fileID, "application/octet-stream", file, func(ctx context.Context, req *http.Request) error {
		// Add saved jwt token for auth.
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	status = resp.StatusCode
	return
}

func GfileList(c *oapi.Client, jwt string) (gfiles map[string]oapi.Gfile, status int, err error) {
	// Create OAPI gfile object.
	resp, err := c.ListGfiles(context.TODO(), func(ctx context.Context, req *http.Request) error {
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

	// Decode gfiles from body.
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&gfiles)
		if err != nil {
			zap.S().Errorln("Can't write to response in Listgfile handler", err)
		}
	}

	return gfiles, resp.StatusCode, nil
}

func GfileGet(c *oapi.Client, conf config.Config, jwt string, fileID, fileName string) (downloaded bool, status int, err error) {
	// Get file.
	resp, err := c.GetGfile(context.TODO(), fileID, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	// Create file
	file, err := os.Create(filepath.Join(conf.FileSavingPath, fileName))
	defer func() {
		if err := file.Close(); err != nil {
			zap.S().Errorln("Can't close file: ", err)
		}
	}()

	if err != nil {
		zap.S().Errorln("Can't create file: ", file)
		return false, http.StatusInternalServerError, err
	}
	// Write file from server.
	n, err := io.Copy(file, resp.Body)
	if err != nil {
		zap.S().Errorln("Can't read from body: ", file)
		return false, http.StatusInternalServerError, err
	}
	zap.S().Infoln("Write file: ", fileName, " ", n)
	return true, resp.StatusCode, nil
}
