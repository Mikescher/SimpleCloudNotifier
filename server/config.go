package server

import (
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	Namespace  string
	GinDebug   bool
	ServerIP   string
	ServerPort string
	DBFile     string
}

var Conf Config

var configLoc = Config{
	Namespace:  "local",
	GinDebug:   true,
	ServerIP:   "0.0.0.0",
	ServerPort: "8080",
	DBFile:     ".run-data/db.sqlite3",
}

var configDev = Config{
	Namespace:  "develop",
	GinDebug:   true,
	ServerIP:   "0.0.0.0",
	ServerPort: "80",
	DBFile:     "/data/scn.sqlite3",
}

var configStag = Config{
	Namespace:  "staging",
	GinDebug:   true,
	ServerIP:   "0.0.0.0",
	ServerPort: "80",
	DBFile:     "/data/scn.sqlite3",
}

var configProd = Config{
	Namespace:  "production",
	GinDebug:   false,
	ServerIP:   "0.0.0.0",
	ServerPort: "80",
	DBFile:     "/data/scn.sqlite3",
}

var allConfig = []Config{
	configLoc,
	configDev,
	configStag,
	configProd,
}

func getConfig(ns string) (Config, bool) {
	if ns == "" {
		return configLoc, true
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
