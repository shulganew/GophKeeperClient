package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/shulganew/GophKeeperClient/internal/app/config"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"go.uber.org/zap"
)

const PreambleLeth = 8
const Content = "application/octet-stream"

// Files add with two steps:
// 1. Uplod file and return created file id in minio storage.
// 2. Create file metadata as sectet in db with users description (definition field and file_id)
type UploadReader struct {
	file      *os.File
	preambule []byte
	metadata  []byte
	index     int64
	metaLen   int64
}

// Constructor for Upload files.
// byte structute: |8-byte preambule with meta length| N-bytes metadata newGfile | File bytes |
func NewUploadReader(file *os.File, preambule []byte, metadata []byte) *UploadReader {
	r := new(UploadReader)
	r.file = file
	r.preambule = preambule
	r.metadata = metadata
	r.metaLen = int64(len(r.metadata))
	return r
}

// Read to b []byte preambule, then metadata, then original file.
func (r *UploadReader) Read(b []byte) (totlal int, err error) {
	// Add preambule bytes (PreambleLeth), witch contains lenth of metadata (newGfile)
	if r.index < PreambleLeth {
		n := copy(b, r.preambule[r.index:PreambleLeth])
		r.index += int64(n)
		totlal += n
	}

	// Add metadata bytes - newGfiles object.
	if r.index >= PreambleLeth && r.index < PreambleLeth+r.metaLen {
		n := copy(b[PreambleLeth:], r.metadata[r.index-PreambleLeth:r.metaLen])
		r.index += int64(n)
		totlal += n
	}
	// Add file bytes
	if r.index >= PreambleLeth+r.metaLen {
		bf := make([]byte, len(b)-totlal)
		_, err := r.file.Read(bf)
		if err != nil {
			return totlal, err
		}
		n := copy(b[PreambleLeth+r.metaLen:], bf)
		r.index += int64(n)
		totlal += n
		return totlal, nil

	}
	return
}

// Upload files to server.
func FileAdd(conf config.Config, jwt, def, fPath string) (status int, err error) {

	// custom HTTP client
	// with a raw http.Response
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(GetTLSClietn()))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(fPath)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Create file
	nfile := oapi.NewGfile{Definition: def, Fname: file.Name()}

	// Encode nfile as metadata to binary
	var md bytes.Buffer
	err = gob.NewEncoder(&md).Encode(&nfile)
	if err != nil {
		return 0, err
	}

	metadata := md.Bytes()
	zap.S().Infoln("metadata", metadata)
	// Write preambule - size of newGfile object (metadata).
	p := make([]byte, PreambleLeth)
	mLen := uint64(len(metadata))
	zap.S().Infoln("Metadata length: ", mLen)
	binary.LittleEndian.PutUint64(p, mLen)

	ur := NewUploadReader(file, p, metadata)
	resp, err := c.AddGfileWithBody(context.TODO(), Content, ur, func(ctx context.Context, req *http.Request) error {
		// Add saved jwt token for auth.
		req.Header.Add("Authorization", config.AuthPrefix+jwt)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Print to log file for debug level.
	for k, v := range resp.Header {
		zap.S().Debugf("%s: %v\r\n", k, v[0])
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Debugln("Body error: ", err.Error())
	}
	zap.S().Debugln("Body: ", string(body))
	zap.S().Debugf("Status Code: %d\r\n", resp.StatusCode)

	// Get JWT token and save to User

	return resp.StatusCode, nil
}

func GfileList(conf config.Config, jwt string) (gfiles []oapi.Gfile, status int, err error) {
	// custom HTTP client
	c, err := oapi.NewClient(conf.Address, oapi.WithHTTPClient(GetTLSClietn()))
	if err != nil {
		log.Fatal(err)
	}

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
