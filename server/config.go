package server

import (
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Config struct {
	Namespace       string
	GinDebug        bool
	ServerIP        string
	ServerPort      string
	DBFile          string
	RequestTimeout  time.Duration
	ReturnRawErrors bool
}

var Conf Config

var configLocHost = Config{
	Namespace:       "local-host",
	GinDebug:        true,
	ServerIP:        "0.0.0.0",
	ServerPort:      "8080",
	DBFile:          ".run-data/db.sqlite3",
	RequestTimeout:  16 * time.Second,
	ReturnRawErrors: true,
}

var configLocDocker = Config{
	Namespace:       "local-docker",
	GinDebug:        true,
	ServerIP:        "0.0.0.0",
	ServerPort:      "80",
	DBFile:          "/data/scn_docker.sqlite3",
	RequestTimeout:  16 * time.Second,
	ReturnRawErrors: true,
}

var configDev = Config{
	Namespace:       "develop",
	GinDebug:        true,
	ServerIP:        "0.0.0.0",
	ServerPort:      "80",
	DBFile:          "/data/scn.sqlite3",
	RequestTimeout:  16 * time.Second,
	ReturnRawErrors: true,
}

var configStag = Config{
	Namespace:       "staging",
	GinDebug:        true,
	ServerIP:        "0.0.0.0",
	ServerPort:      "80",
	DBFile:          "/data/scn.sqlite3",
	RequestTimeout:  16 * time.Second,
	ReturnRawErrors: true,
}

var configProd = Config{
	Namespace:       "production",
	GinDebug:        false,
	ServerIP:        "0.0.0.0",
	ServerPort:      "80",
	DBFile:          "/data/scn.sqlite3",
	RequestTimeout:  16 * time.Second,
	ReturnRawErrors: false,
}

var allConfig = []Config{
	configLocHost,
	configLocDocker,
	configDev,
	configStag,
	configProd,
}

func getConfig(ns string) (Config, bool) {
	if ns == "" {
		return configLocHost, true
	}
	for _, c := range allConfig {
		if c.Namespace == ns {
			return c, true
		}
	}
	return Config{}, false
}

func init() {
	ns := os.Getenv("CONF_NS")

	cfg, ok := getConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}

	Conf = cfg
}
