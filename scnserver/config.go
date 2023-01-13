package server

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/confext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"os"
	"time"
)

type Config struct {
	Namespace                string
	BaseURL                  string        `env:"SCN_URL"`
	GinDebug                 bool          `env:"SCN_GINDEBUG"`
	LogLevel                 zerolog.Level `env:"SCN_LOGLEVEL"`
	ServerIP                 string        `env:"SCN_IP"`
	ServerPort               string        `env:"SCN_PORT"`
	DBMain                   DBConfig      `env:"SCN_DB_MAIN"`
	DBRequests               DBConfig      `env:"SCN_DB_REQUESTS"`
	DBLogs                   DBConfig      `env:"SCN_DB_LOGS"`
	RequestTimeout           time.Duration `env:"SCN_REQUEST_TIMEOUT"`
	RequestMaxRetry          int           `env:"SCN_REQUEST_MAXRETRY"`
	RequestRetrySleep        time.Duration `env:"SCN_REQUEST_RETRYSLEEP"`
	Cors                     bool          `env:"SCN_CORS"`
	ReturnRawErrors          bool          `env:"SCN_ERROR_RETURN"`
	DummyFirebase            bool          `env:"SCN_DUMMY_FB"`
	DummyGoogleAPI           bool          `env:"SCN_DUMMY_GOOG"`
	FirebaseTokenURI         string        `env:"SCN_FB_TOKENURI"`
	FirebaseProjectID        string        `env:"SCN_FB_PROJECTID"`
	FirebasePrivKeyID        string        `env:"SCN_FB_PRIVATEKEYID"`
	FirebaseClientMail       string        `env:"SCN_FB_CLIENTEMAIL"`
	FirebasePrivateKey       string        `env:"SCN_FB_PRIVATEKEY"`
	GoogleAPITokenURI        string        `env:"SCN_GOOG_TOKENURI"`
	GoogleAPIPrivKeyID       string        `env:"SCN_GOOG_PRIVATEKEYID"`
	GoogleAPIClientMail      string        `env:"SCN_GOOG_CLIENTEMAIL"`
	GoogleAPIPrivateKey      string        `env:"SCN_GOOG_PRIVATEKEY"`
	GooglePackageName        string        `env:"SCN_GOOG_PACKAGENAME"`
	GoogleProProductID       string        `env:"SCN_GOOG_PROPRODUCTID"`
	ReqLogEnabled            bool          `env:"SCN_REQUESTLOG_ENABLED"`
	ReqLogMaxBodySize        int           `env:"SCN_REQUESTLOG_MAXBODYSIZE"`
	ReqLogHistoryMaxCount    int           `env:"SCN_REQUESTLOG_HISTORY_MAXCOUNT"`
	ReqLogHistoryMaxDuration time.Duration `env:"SCN_REQUESTLOG_HISTORY_MAXDURATION"`
}

type DBConfig struct {
	File             string        `env:"FILE"`
	Journal          string        `env:"JOURNAL"`
	Timeout          time.Duration `env:"TIMEOUT"`
	MaxOpenConns     int           `env:"MAXOPENCONNECTIONS"`
	MaxIdleConns     int           `env:"MAXIDLECONNECTIONS"`
	ConnMaxLifetime  time.Duration `env:"CONNEXTIONMAXLIFETIME"`
	ConnMaxIdleTime  time.Duration `env:"CONNEXTIONMAXIDLETIME"`
	CheckForeignKeys bool          `env:"CHECKFOREIGNKEYS"`
	SingleConn       bool          `env:"SINGLECONNECTION"`
}

var Conf Config

var configLocHost = func() Config {
	return Config{
		Namespace:  "local-host",
		BaseURL:    "http://localhost:8080",
		GinDebug:   false,
		LogLevel:   zerolog.DebugLevel,
		ServerIP:   "0.0.0.0",
		ServerPort: "8080",
		DBMain: DBConfig{
			File:             ".run-data/loc_main.sqlite3",
			Journal:          "WAL",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBRequests: DBConfig{
			File:             ".run-data/loc_requests.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBLogs: DBConfig{
			File:             ".run-data/loc_logs.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          true,
		DummyFirebase:            true,
		FirebaseTokenURI:         "",
		FirebaseProjectID:        "",
		FirebasePrivKeyID:        "",
		FirebaseClientMail:       "",
		FirebasePrivateKey:       "",
		DummyGoogleAPI:           true,
		GoogleAPITokenURI:        "",
		GoogleAPIPrivKeyID:       "",
		GoogleAPIClientMail:      "",
		GoogleAPIPrivateKey:      "",
		GooglePackageName:        "",
		GoogleProProductID:       "",
		Cors:                     true,
		ReqLogEnabled:            true,
		ReqLogMaxBodySize:        2048,
		ReqLogHistoryMaxCount:    1638,
		ReqLogHistoryMaxDuration: timeext.FromDays(60),
	}
}

var configLocDocker = func() Config {
	return Config{
		Namespace:  "local-docker",
		BaseURL:    "http://localhost:8080",
		GinDebug:   false,
		LogLevel:   zerolog.DebugLevel,
		ServerIP:   "0.0.0.0",
		ServerPort: "80",
		DBMain: DBConfig{
			File:             "/data/docker_scn_main.sqlite3",
			Journal:          "WAL",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBRequests: DBConfig{
			File:             "/data/docker_scn_requests.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBLogs: DBConfig{
			File:             "/data/docker_scn_logs.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          true,
		DummyFirebase:            true,
		FirebaseTokenURI:         "",
		FirebaseProjectID:        "",
		FirebasePrivKeyID:        "",
		FirebaseClientMail:       "",
		FirebasePrivateKey:       "",
		DummyGoogleAPI:           true,
		GoogleAPITokenURI:        "",
		GoogleAPIPrivKeyID:       "",
		GoogleAPIClientMail:      "",
		GoogleAPIPrivateKey:      "",
		GooglePackageName:        "",
		GoogleProProductID:       "",
		Cors:                     true,
		ReqLogMaxBodySize:        2048,
		ReqLogHistoryMaxCount:    1638,
		ReqLogHistoryMaxDuration: timeext.FromDays(60),
	}
}

var configDev = func() Config {
	return Config{
		Namespace:  "develop",
		BaseURL:    confEnv("SCN_URL"),
		GinDebug:   false,
		LogLevel:   zerolog.DebugLevel,
		ServerIP:   "0.0.0.0",
		ServerPort: "80",
		DBMain: DBConfig{
			File:             "/data/scn_main.sqlite3",
			Journal:          "WAL",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBRequests: DBConfig{
			File:             "/data/scn_requests.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBLogs: DBConfig{
			File:             "/data/scn_logs.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          true,
		DummyFirebase:            false,
		FirebaseTokenURI:         "https://oauth2.googleapis.com/token",
		FirebaseProjectID:        confEnv("SCN_FB_PROJECTID"),
		FirebasePrivKeyID:        confEnv("SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:       confEnv("SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:       confEnv("SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:           false,
		GoogleAPITokenURI:        "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:       confEnv("SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail:      confEnv("SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey:      confEnv("SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:        confEnv("SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:       confEnv("SCN_GOOG_PROPRODUCTID"),
		Cors:                     true,
		ReqLogMaxBodySize:        2048,
		ReqLogHistoryMaxCount:    1638,
		ReqLogHistoryMaxDuration: timeext.FromDays(60),
	}
}

var configStag = func() Config {
	return Config{
		Namespace:  "staging",
		BaseURL:    confEnv("SCN_URL"),
		GinDebug:   false,
		LogLevel:   zerolog.DebugLevel,
		ServerIP:   "0.0.0.0",
		ServerPort: "80",
		DBMain: DBConfig{
			File:             "/data/scn_main.sqlite3",
			Journal:          "WAL",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBRequests: DBConfig{
			File:             "/data/scn_requests.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBLogs: DBConfig{
			File:             "/data/scn_logs.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          true,
		DummyFirebase:            false,
		FirebaseTokenURI:         "https://oauth2.googleapis.com/token",
		FirebaseProjectID:        confEnv("SCN_FB_PROJECTID"),
		FirebasePrivKeyID:        confEnv("SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:       confEnv("SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:       confEnv("SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:           false,
		GoogleAPITokenURI:        "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:       confEnv("SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail:      confEnv("SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey:      confEnv("SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:        confEnv("SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:       confEnv("SCN_GOOG_PROPRODUCTID"),
		Cors:                     true,
		ReqLogMaxBodySize:        2048,
		ReqLogHistoryMaxCount:    1638,
		ReqLogHistoryMaxDuration: timeext.FromDays(60),
	}
}

var configProd = func() Config {
	return Config{
		Namespace:  "production",
		BaseURL:    confEnv("SCN_URL"),
		GinDebug:   false,
		LogLevel:   zerolog.InfoLevel,
		ServerIP:   "0.0.0.0",
		ServerPort: "80",
		DBMain: DBConfig{
			File:             "/data/scn_main.sqlite3",
			Journal:          "WAL",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBRequests: DBConfig{
			File:             "/data/scn_requests.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		DBLogs: DBConfig{
			File:             "/data/scn_logs.sqlite3",
			Journal:          "DELETE",
			Timeout:          5 * time.Second,
			CheckForeignKeys: false,
			SingleConn:       false,
			MaxOpenConns:     5,
			MaxIdleConns:     5,
			ConnMaxLifetime:  60 * time.Minute,
			ConnMaxIdleTime:  60 * time.Minute,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          false,
		DummyFirebase:            false,
		FirebaseTokenURI:         "https://oauth2.googleapis.com/token",
		FirebaseProjectID:        confEnv("SCN_SCN_FB_PROJECTID"),
		FirebasePrivKeyID:        confEnv("SCN_SCN_FB_PRIVATEKEYID"),
		FirebaseClientMail:       confEnv("SCN_SCN_FB_CLIENTEMAIL"),
		FirebasePrivateKey:       confEnv("SCN_SCN_FB_PRIVATEKEY"),
		DummyGoogleAPI:           false,
		GoogleAPITokenURI:        "https://oauth2.googleapis.com/token",
		GoogleAPIPrivKeyID:       confEnv("SCN_SCN_GOOG_PRIVATEKEYID"),
		GoogleAPIClientMail:      confEnv("SCN_SCN_GOOG_CLIENTEMAIL"),
		GoogleAPIPrivateKey:      confEnv("SCN_SCN_GOOG_PRIVATEKEY"),
		GooglePackageName:        confEnv("SCN_SCN_GOOG_PACKAGENAME"),
		GoogleProProductID:       confEnv("SCN_SCN_GOOG_PROPRODUCTID"),
		Cors:                     true,
		ReqLogMaxBodySize:        2048,
		ReqLogHistoryMaxCount:    1638,
		ReqLogHistoryMaxDuration: timeext.FromDays(60),
	}
}

var allConfig = map[string]func() Config{
	"local-host":   configLocHost,
	"local-docker": configLocDocker,
	"develop":      configDev,
	"staging":      configStag,
	"production":   configProd,
}

func GetConfig(ns string) (Config, bool) {
	if ns == "" {
		ns = "local-host"
	}
	if cfn, ok := allConfig[ns]; ok {
		c := cfn()
		err := confext.ApplyEnvOverrides(&c, "_")
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

	cfg, ok := GetConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}

	Conf = cfg
}
