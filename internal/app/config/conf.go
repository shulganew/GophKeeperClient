package config

import (
	"flag"
	"net/url"
	"time"

	"github.com/shulganew/GophKeeperClient/internal/app/validators"
	"go.uber.org/zap"
)

const AuthPrefix = "Bearer "
const Shema = "https"
const TokenExp = time.Hour * 3600
const DataBaseType = "postgres"

type Config struct {
	// flag -a, Server address
	Address        string
	FileSavingPath string
	SertPath       string //sertificate TLS file path (server public key)
}

func InitConfig() *Config {
	config := Config{}
	// read command line argue
	serverAddress := flag.String("a", "localhost:8443", "Service GKeeper address")
	filePath := flag.String("f", "/home/igor/files/", "Service GKeeper address")
	sertPath := flag.String("s", "cert/server.crt", "Service GKeeper address")
	flag.Parse()

	// Check and parse URL
	startaddr, startport := validators.CheckURL(*serverAddress)
	// Server address
	u := url.URL{Scheme: Shema, Host: startaddr + ":" + startport}
	config.Address = u.String()
	config.FileSavingPath = *filePath
	config.SertPath = *sertPath
	zap.S().Infoln("Configuration complite")
	return &config
}
