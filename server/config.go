package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Config struct {
	Namespace          string
	GinDebug           bool
	ServerIP           string
	ServerPort         string
	DBFile             string
	RequestTimeout     time.Duration
	ReturnRawErrors    bool
	DummyFirebase      bool
	FirebaseTokenURI   string
	FirebaseProjectID  string
	FirebasePrivKeyID  string
	FirebaseClientMail string
	FirebasePrivateKey string
}

var Conf Config

var configLocHost = func() Config {
	return Config{
		Namespace:          "local-host",
		GinDebug:           true,
		ServerIP:           "0.0.0.0",
		ServerPort:         "8080",
		DBFile:             ".run-data/db.sqlite3",
		RequestTimeout:     16 * time.Second,
		ReturnRawErrors:    true,
		DummyFirebase:      true,
		FirebaseTokenURI:   "",
		FirebaseProjectID:  "",
		FirebasePrivKeyID:  "",
		FirebaseClientMail: "",
		FirebasePrivateKey: "",
	}
}

var configLocDocker = func() Config {
	return Config{
		Namespace:          "local-docker",
		GinDebug:           true,
		ServerIP:           "0.0.0.0",
		ServerPort:         "80",
		DBFile:             "/data/scn_docker.sqlite3",
		RequestTimeout:     16 * time.Second,
		ReturnRawErrors:    true,
		DummyFirebase:      true,
		FirebaseTokenURI:   "",
		FirebaseProjectID:  "",
		FirebasePrivKeyID:  "",
		FirebaseClientMail: "",
		FirebasePrivateKey: "",
	}
}

var configDev = func() Config {
	return Config{
		Namespace:          "develop",
		GinDebug:           true,
		ServerIP:           "0.0.0.0",
		ServerPort:         "80",
		DBFile:             "/data/scn.sqlite3",
		RequestTimeout:     16 * time.Second,
		ReturnRawErrors:    true,
		DummyFirebase:      false,
		FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
		FirebaseProjectID:  confEnv("FB_PROJECTID"),
		FirebasePrivKeyID:  confEnv("FB_PRIVATEKEYID"),
		FirebaseClientMail: confEnv("FB_CLIENTEMAIL"),
		FirebasePrivateKey: confEnv("FB_PRIVATEKEY"),
	}
}

var configStag = func() Config {
	return Config{
		Namespace:          "staging",
		GinDebug:           true,
		ServerIP:           "0.0.0.0",
		ServerPort:         "80",
		DBFile:             "/data/scn.sqlite3",
		RequestTimeout:     16 * time.Second,
		ReturnRawErrors:    true,
		DummyFirebase:      false,
		FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
		FirebaseProjectID:  confEnv("FB_PROJECTID"),
		FirebasePrivKeyID:  confEnv("FB_PRIVATEKEYID"),
		FirebaseClientMail: confEnv("FB_CLIENTEMAIL"),
		FirebasePrivateKey: confEnv("FB_PRIVATEKEY"),
	}
}

var configProd = func() Config {
	return Config{
		Namespace:          "production",
		GinDebug:           false,
		ServerIP:           "0.0.0.0",
		ServerPort:         "80",
		DBFile:             "/data/scn.sqlite3",
		RequestTimeout:     16 * time.Second,
		ReturnRawErrors:    false,
		DummyFirebase:      false,
		FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
		FirebaseProjectID:  confEnv("FB_PROJECTID"),
		FirebasePrivKeyID:  confEnv("FB_PRIVATEKEYID"),
		FirebaseClientMail: confEnv("FB_CLIENTEMAIL"),
		FirebasePrivateKey: confEnv("FB_PRIVATEKEY"),
	}
}

var allConfig = map[string]func() Config{
	"local-host":   configLocHost,
	"local-docker": configLocDocker,
	"develop":      configDev,
	"staging":      configStag,
	"production":   configProd,
}

func getConfig(ns string) (Config, bool) {
	if ns == "" {
		return configLocHost(), true
	}
	if c, ok := allConfig[ns]; ok {
		return c(), true
	}
	return Config{}, false
}

func confEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	} else {
		panic(fmt.Sprintf("Missing required environment variable '%s'", key))
	}
}

func init() {
	ns := os.Getenv("CONF_NS")

	cfg, ok := getConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}

	Conf = cfg
}
