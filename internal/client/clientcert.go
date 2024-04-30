package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func GetTLSClietn() *http.Client {
	caCertf, _ := os.ReadFile("./cert/server.crt")
	rootCAs, _ := x509.SystemCertPool()
	// handle case where rootCAs == nil and create an empty pool...
	if ok := rootCAs.AppendCertsFromPEM(caCertf); !ok {
		zap.S().Infoln("Can't load trasted sertifecate!")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            rootCAs,
	}

	hc := &http.Client{
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}
	return hc
}
