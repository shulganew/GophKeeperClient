package config

import (
	"flag"
	"os"
	"time"

	"github.com/shulganew/GophKeeperClient/internal/app/validators"
	"go.uber.org/zap"
)

const RegisterPath = "/api/user/register"
const LoginPath = "/api/user/login"
const Shema = "http"
const TokenExp = time.Hour * 3600
const DataBaseType = "postgres"

type Config struct {
	// flag -a, Server address
	Address string
}

func InitConfig() *Config {
	config := Config{}
	// read command line argue
	serverAddress := flag.String("a", "localhost:8080", "Service GKeeper address")
	flag.Parse()

	// Check and parse URL
	startaddr, startport := validators.CheckURL(*serverAddress)
	// Server address
	config.Address = startaddr + ":" + startport

	// read OS ENVs
	addr, exist := os.LookupEnv(("RUN_ADDRESS"))

	// if env var does not exist  - set def value
	if exist {
		config.Address = addr
		zap.S().Infoln("Set result address from evn RUN_ADDRESS: ", config.Address)
	} else {
		zap.S().Infoln("Env var RUN_ADDRESS not found, use default", config.Address)
	}

	zap.S().Infoln("Configuration complite")
	return &config
}
