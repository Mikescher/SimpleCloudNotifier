package server

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/confext"
	"os"
	"time"
)

type Config struct {
	Namespace           string
	BaseURL             string        `env:"SCN_URL"`
	GinDebug            bool          `env:"SCN_GINDEBUG"`
	LogLevel            zerolog.Level `env:"SCN_LOGLEVEL"`
	ServerIP            string        `env:"SCN_IP"`
	ServerPort          string        `env:"SCN_PORT"`
	DBFile              string        `env:"SCN_DB_FILE"`
	DBJournal           string        `env:"SCN_DB_JOURNAL"`
	DBTimeout           time.Duration `env:"SCN_DB_TIMEOUT"`
	DBMaxOpenConns      int           `env:"SCN_DB_MAXOPENCONNECTIONS"`
	DBMaxIdleConns      int           `env:"SCN_DB_MAXIDLECONNECTIONS"`
	DBConnMaxLifetime   time.Duration `env:"SCN_DB_CONNEXTIONMAXLIFETIME"`
	DBConnMaxIdleTime   time.Duration `env:"SCN_DB_CONNEXTIONMAXIDLETIME"`
	DBCheckForeignKeys  bool          `env:"SCN_DB_CHECKFOREIGNKEYS"`
	DBSingleConn        bool          `env:"SCN_DB_SINGLECONNECTION"`
	RequestTimeout      time.Duration `env:"SCN_REQUEST_TIMEOUT"`
	RequestMaxRetry     int           `env:"SCN_REQUEST_MAXRETRY"`
	RequestRetrySleep   time.Duration `env:"SCN_REQUEST_RETRYSLEEP"`
	ReturnRawErrors     bool          `env:"SCN_ERROR_RETURN"`
	DummyFirebase       bool          `env:"SCN_DUMMY_FB"`
	DummyGoogleAPI      bool          `env:"SCN_DUMMY_GOOG"`
	FirebaseTokenURI    string        `env:"SCN_FB_TOKENURI"`
	FirebaseProjectID   string        `env:"SCN_FB_PROJECTID"`
	FirebasePrivKeyID   string        `env:"SCN_FB_PRIVATEKEYID"`
	FirebaseClientMail  string        `env:"SCN_FB_CLIENTEMAIL"`
	FirebasePrivateKey  string        `env:"SCN_FB_PRIVATEKEY"`
	GoogleAPITokenURI   string        `env:"SCN_GOOG_TOKENURI"`
	GoogleAPIPrivKeyID  string        `env:"SCN_GOOG_PRIVATEKEYID"`
	GoogleAPIClientMail string        `env:"SCN_GOOG_CLIENTEMAIL"`
	GoogleAPIPrivateKey string        `env:"SCN_GOOG_PRIVATEKEY"`
	GooglePackageName   string        `env:"SCN_GOOG_PACKAGENAME"`
	GoogleProProductID  string        `env:"SCN_GOOG_PROPRODUCTID"`
	Cors                bool          `env:"SCN_CORS"`
}

var Conf Config

var configLocHost = func() Config {
	return Config{
		Namespace:           "local-host",
		BaseURL:             "http://localhost:8080",
		GinDebug:            true,
		LogLevel:            zerolog.DebugLevel,
		ServerIP:            "0.0.0.0",
		ServerPort:          "8080",
		DBFile:              ".run-data/db.sqlite3",
		DBJournal:           "WAL",
		DBTimeout:           5 * time.Second,
		DBCheckForeignKeys:  false,
		DBSingleConn:        false,
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      5,
		DBConnMaxLifetime:   60 * time.Minute,
		DBConnMaxIdleTime:   60 * time.Minute,
		RequestTimeout:      16 * time.Second,
		RequestMaxRetry:     8,
		RequestRetrySleep:   100 * time.Millisecond,
		ReturnRawErrors:     true,
		DummyFirebase:       true,
		FirebaseTokenURI:    "",
		FirebaseProjectID:   "",
		FirebasePrivKeyID:   "",
		FirebaseClientMail:  "",
		FirebasePrivateKey:  "",
		DummyGoogleAPI:      true,
		GoogleAPITokenURI:   "",
		GoogleAPIPrivKeyID:  "",
		GoogleAPIClientMail: "",
		GoogleAPIPrivateKey: "",
		GooglePackageName:   "",
		GoogleProProductID:  "",
		Cors:                true,
	}
}

var configLocDocker = func() Config {
	return Config{
		Namespace:           "local-docker",
		BaseURL:             "http://localhost:8080",
		GinDebug:            true,
		LogLevel:            zerolog.DebugLevel,
		ServerIP:            "0.0.0.0",
		ServerPort:          "80",
		DBFile:              "/data/scn_docker.sqlite3",
		DBJournal:           "WAL",
		DBTimeout:           5 * time.Second,
		DBCheckForeignKeys:  false,
		DBSingleConn:        false,
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      5,
		DBConnMaxLifetime:   60 * time.Minute,
		DBConnMaxIdleTime:   60 * time.Minute,
		RequestTimeout:      16 * time.Second,
		RequestMaxRetry:     8,
		RequestRetrySleep:   100 * time.Millisecond,
		ReturnRawErrors:     true,
		DummyFirebase:       true,
		FirebaseTokenURI:    "",
		FirebaseProjectID:   "",
		FirebasePrivKeyID:   "",
		FirebaseClientMail:  "",
		FirebasePrivateKey:  "",
		DummyGoogleAPI:      true,
		GoogleAPITokenURI:   "",
		GoogleAPIPrivKeyID:  "",
		GoogleAPIClientMail: "",
		GoogleAPIPrivateKey: "",
		GooglePackageName:   "",
		GoogleProProductID:  "",
		Cors:                true,
	}
}

var configDev = func() Config {
	return Config{
		Namespace:           "develop",
		BaseURL:             confEnv("SCN_URL"),
		GinDebug:            true,
		LogLevel:            zerolog.DebugLevel,
		ServerIP:            "0.0.0.0",
		ServerPort:          "80",
		DBFile:              "/data/scn.sqlite3",
		DBJournal:           "WAL",
		DBTimeout:           5 * time.Second,
		DBCheckForeignKeys:  false,
		DBSingleConn:        false,
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      5,
		DBConnMaxLifetime:   60 * time.Minute,
		DBConnMaxIdleTime:   60 * time.Minute,
		RequestTimeout:      16 * time.Second,
		RequestMaxRetry:     8,
		RequestRetrySleep:   100 * time.Millisecond,
		ReturnRawErrors:     true,
		DummyFirebase:       false,
		FirebaseTokenURI:    "https://oauth2.googleapis.com/token",
		FirebaseProjectID:   confEnv("SCN_FB_PROJECTID"),
		FirebasePrivKeyID:   confEnv("SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:  confEnv("SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:  confEnv("SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:      false,
		GoogleAPITokenURI:   "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:  confEnv("SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail: confEnv("SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey: confEnv("SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:   confEnv("SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:  confEnv("SCN_GOOG_PROPRODUCTID"),
		Cors:                true,
	}
}

var configStag = func() Config {
	return Config{
		Namespace:           "staging",
		BaseURL:             confEnv("SCN_URL"),
		GinDebug:            true,
		LogLevel:            zerolog.DebugLevel,
		ServerIP:            "0.0.0.0",
		ServerPort:          "80",
		DBFile:              "/data/scn.sqlite3",
		DBJournal:           "WAL",
		DBTimeout:           5 * time.Second,
		DBCheckForeignKeys:  false,
		DBSingleConn:        false,
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      5,
		DBConnMaxLifetime:   60 * time.Minute,
		DBConnMaxIdleTime:   60 * time.Minute,
		RequestTimeout:      16 * time.Second,
		RequestMaxRetry:     8,
		RequestRetrySleep:   100 * time.Millisecond,
		ReturnRawErrors:     true,
		DummyFirebase:       false,
		FirebaseTokenURI:    "https://oauth2.googleapis.com/token",
		FirebaseProjectID:   confEnv("SCN_FB_PROJECTID"),
		FirebasePrivKeyID:   confEnv("SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:  confEnv("SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:  confEnv("SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:      false,
		GoogleAPITokenURI:   "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:  confEnv("SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail: confEnv("SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey: confEnv("SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:   confEnv("SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:  confEnv("SCN_GOOG_PROPRODUCTID"),
		Cors:                true,
	}
}

var configProd = func() Config {
	return Config{
		Namespace:           "production",
		BaseURL:             confEnv("SCN_URL"),
		GinDebug:            false,
		LogLevel:            zerolog.InfoLevel,
		ServerIP:            "0.0.0.0",
		ServerPort:          "80",
		DBFile:              "/data/scn.sqlite3",
		DBJournal:           "WAL",
		DBTimeout:           5 * time.Second,
		DBCheckForeignKeys:  false,
		DBSingleConn:        false,
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      5,
		DBConnMaxLifetime:   60 * time.Minute,
		DBConnMaxIdleTime:   60 * time.Minute,
		RequestTimeout:      16 * time.Second,
		RequestMaxRetry:     8,
		RequestRetrySleep:   100 * time.Millisecond,
		ReturnRawErrors:     false,
		DummyFirebase:       false,
		FirebaseTokenURI:    "https://oauth2.googleapis.com/token",
		FirebaseProjectID:   confEnv("SCN_SCN_FB_PROJECTID"),
		FirebasePrivKeyID:   confEnv("SCN_SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:  confEnv("SCN_SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:  confEnv("SCN_SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:      false,
		GoogleAPITokenURI:   "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:  confEnv("SCN_SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail: confEnv("SCN_SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey: confEnv("SCN_SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:   confEnv("SCN_SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:  confEnv("SCN_SCN_GOOG_PROPRODUCTID"),
		Cors:                true,
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
		ns = "local-host"
	}
	if cfn, ok := allConfig[ns]; ok {
		c := cfn()
		err := confext.ApplyEnvOverrides(&c)
		if err != nil {
			panic(err)
		}
		return c, true
	}
	return Config{}, false
}

func confEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	} else {
		log.Fatal().Msg(fmt.Sprintf("Missing required environment variable '%s'", key))
		return ""
	}
}

func init() {
	ns := os.Getenv("SCN_NAMESPACE")

	cfg, ok := getConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}

	Conf = cfg
}
