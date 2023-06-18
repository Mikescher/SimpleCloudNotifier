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
	BaseURL                  string        `env:"URL"`
	GinDebug                 bool          `env:"GINDEBUG"`
	LogLevel                 zerolog.Level `env:"LOGLEVEL"`
	ServerIP                 string        `env:"IP"`
	ServerPort               string        `env:"PORT"`
	DBMain                   DBConfig      `env:"DB_MAIN"`
	DBRequests               DBConfig      `env:"DB_REQUESTS"`
	DBLogs                   DBConfig      `env:"DB_LOGS"`
	RequestTimeout           time.Duration `env:"REQUEST_TIMEOUT"`
	RequestMaxRetry          int           `env:"REQUEST_MAXRETRY"`
	RequestRetrySleep        time.Duration `env:"REQUEST_RETRYSLEEP"`
	Cors                     bool          `env:"CORS"`
	ReturnRawErrors          bool          `env:"ERROR_RETURN"`
	DummyFirebase            bool          `env:"DUMMY_FB"`
	DummyGoogleAPI           bool          `env:"DUMMY_GOOG"`
	FirebaseTokenURI         string        `env:"FB_TOKENURI"`
	FirebaseProjectID        string        `env:"FB_PROJECTID"`
	FirebasePrivKeyID        string        `env:"FB_PRIVATEKEYID"`
	FirebaseClientMail       string        `env:"FB_CLIENTEMAIL"`
	FirebasePrivateKey       string        `env:"FB_PRIVATEKEY"`
	GoogleAPITokenURI        string        `env:"GOOG_TOKENURI"`
	GoogleAPIPrivKeyID       string        `env:"GOOG_PRIVATEKEYID"`
	GoogleAPIClientMail      string        `env:"GOOG_CLIENTEMAIL"`
	GoogleAPIPrivateKey      string        `env:"GOOG_PRIVATEKEY"`
	GooglePackageName        string        `env:"GOOG_PACKAGENAME"`
	GoogleProProductID       string        `env:"GOOG_PROPRODUCTID"`
	ReqLogEnabled            bool          `env:"REQUESTLOG_ENABLED"`
	ReqLogMaxBodySize        int           `env:"REQUESTLOG_MAXBODYSIZE"`
	ReqLogHistoryMaxCount    int           `env:"REQUESTLOG_HISTORY_MAXCOUNT"`
	ReqLogHistoryMaxDuration time.Duration `env:"REQUESTLOG_HISTORY_MAXDURATION"`
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
	BusyTimeout      time.Duration `env:"BUSYTIMEOUT"`
	EnableLogger     bool          `env:"ENABLELOGGER"`
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
			BusyTimeout:      100 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      100 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      100 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
		ReqLogEnabled:            true,
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
			BusyTimeout:      100 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
		ReqLogEnabled:            true,
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
			BusyTimeout:      100 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
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
			BusyTimeout:      500 * time.Millisecond,
			EnableLogger:     true,
		},
		RequestTimeout:           16 * time.Second,
		RequestMaxRetry:          8,
		RequestRetrySleep:        100 * time.Millisecond,
		ReturnRawErrors:          false,
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
		ReqLogEnabled:            true,
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
		err := confext.ApplyEnvOverrides("SCN_", &c, "_")
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
	ns := os.Getenv("CONF_NS")

	cfg, ok := GetConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}

	Conf = cfg
}
