package server

import (
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
	FirebaseProjectID  string
	FirebaseTokenURI   string
	FirebasePrivKeyID  string
	FirebaseClientMail string
	FirebasePrivateKey string
}

var Conf Config

var configLocHost = Config{
	Namespace:          "local-host",
	GinDebug:           true,
	ServerIP:           "0.0.0.0",
	ServerPort:         "8080",
	DBFile:             ".run-data/db.sqlite3",
	RequestTimeout:     16 * time.Second,
	ReturnRawErrors:    true,
	FirebaseProjectID:  "simplecloudnotifier-ea7ef",
	FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
	FirebasePrivKeyID:  "5bfab19fca25034e87c5b3bd1a4334499d2d1f85",
	FirebaseClientMail: "firebase-adminsdk-42grv@simplecloudnotifier-ea7ef.iam.gserviceaccount.com",
	FirebasePrivateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQD2NWOQDcalRdkp\nHtQHABLlu3GMBQBJrGiCxzOZhi/lLwrw2MJEmg1VFz6TVkX2z3SCzXCPOgGriM70\nuWCNLyZQvUng7u6/WH9hlpCg0vJpkw6BvOBt1zYu3gbb5M0SKEOR+lDVccEjAnT4\nexebXdJHJcbaYAcPnBQ9tgP+cozQBnr2EfxYL0bGMgiH9fErJSGMBDFI996uUW9a\nbtfkZ/XpZqYAvyGQMEjknGnQ8t8PHAnsS9dc1PXSWfBvz07ba3fkypWcpTsIYUiZ\nSpwTLV8awihKHJuphoTWb4x6p/ijop05qr1p3fe8gZd9qOGgALe+JT4IBLgNYKrP\nLMSKH3TdAgMBAAECggEAdFcWDOP1kfNHgl7G4efvBg9kwD08vZNybxmiEFGQIEPy\nb4x9f90rn6G0N/r0ZIPzEjvxjDxkvaGP6aQPM6er+0r2tgsxVcmDp6F2Bgin86tB\nl5ygkEa5m7vekdmz7XiJNVmLCNEP6nMmwqOnrArRaj03kcj+jSm7hs2TZZDLaSA5\nf+2q7h0jaU7Nm0ZwCNJqfPJEGdu1J3fR29Ej0rI8N0w/BuYRet1VYDO09lquqOPS\n0WirOOWV6eyqijqRT+RCt0vVzAppS6guhN7J7RS0V9GLJ/13sdvHuJy/WTjBb7gQ\na6QTo8D3yYF+cn3+0BmgP55uW7N6tsYwXIRZcTI3IQKBgQD+tDKMx0puZu+8zTX9\nC2oHSb4Frl2xq17ZpbkfFmOBPWfQbAHNiQTUoQlzCOQM6QejykXFvfsddP7EY2tL\npgLUrBh81wSCAOOo19vYwQB3YKa5ZZucKxh2VxFSefL/+BYHijFb0mWBj5HmqWS6\n7l6IYT3L04aRK9kxj0Cg6L/z6wKBgQD3dh/kQlPemfdxRpZUJ6WEE5x3Bv7WjLop\nnWgE02Pk8+DB+s50GD3nOR276ADCYS6OkBsgfMkwhhKWZigiEoK9DMul5n587jc9\no5AalZN3IbBGAoXk+u3g1GC9bOY3454K6IJyhehDTImEFyfm00qfUL8fMNcdEx8O\nnwxtyRawVwKBgGqsnd9IOGw0wIOajtoERcv3npZSiPs4guk092uFvPcL+MbZ9YdX\ns6Y6K/L57klZ79ExjjdbcijML0ehO/ba+KSJz1e51jF8ndzBS1pkuwVEfY94dsvZ\nYM1vednJKXT7On696h5C6DBzKPAqUf3Yh88mqvMLDHkQnE6daLv7vykxAoGAOPmA\ndDx1NO48E1+OIwgRyqv9PUZmDB3Qit5L4biN6lvgJqlJOV+PeRokZ2wOKLLZVkeF\nh2BTrhFgXDJfESEz6rT0eljsTHVIUK/E8On5Ttd5z1SrYUII3NfpAhP9mWaVr6tC\nxX1hMYWAr+Ho9PM23iFoL5U+IdqSLvqdkPVYfPcCgYB1ANKNYPIJNx/wLxYWNS0r\nI98HwKfv2TxxE/l+2459NMMHY5wlpFl7MNoeK2SdY+ghWPlxC6u5Nxpnk+bZ8TJe\np7U2nY0SQDLCmPgGWs3KBb/zR49X2b7JS3CXXqQSrLxBe2phZg6kE5nB6NPUDc/i\n6WG8tG20rCfgwlXeXl0+Ow==\n-----END PRIVATE KEY-----\n",
}

var configLocDocker = Config{
	Namespace:          "local-docker",
	GinDebug:           true,
	ServerIP:           "0.0.0.0",
	ServerPort:         "80",
	DBFile:             "/data/scn_docker.sqlite3",
	RequestTimeout:     16 * time.Second,
	ReturnRawErrors:    true,
	FirebaseProjectID:  "simplecloudnotifier-ea7ef",
	FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
	FirebasePrivKeyID:  "5bfab19fca25034e87c5b3bd1a4334499d2d1f85",
	FirebaseClientMail: "firebase-adminsdk-42grv@simplecloudnotifier-ea7ef.iam.gserviceaccount.com",
	FirebasePrivateKey: "TODO",
}

var configDev = Config{
	Namespace:          "develop",
	GinDebug:           true,
	ServerIP:           "0.0.0.0",
	ServerPort:         "80",
	DBFile:             "/data/scn.sqlite3",
	RequestTimeout:     16 * time.Second,
	ReturnRawErrors:    true,
	FirebaseProjectID:  "simplecloudnotifier-ea7ef",
	FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
	FirebasePrivKeyID:  "5bfab19fca25034e87c5b3bd1a4334499d2d1f85",
	FirebaseClientMail: "firebase-adminsdk-42grv@simplecloudnotifier-ea7ef.iam.gserviceaccount.com",
	FirebasePrivateKey: "TODO",
}

var configStag = Config{
	Namespace:          "staging",
	GinDebug:           true,
	ServerIP:           "0.0.0.0",
	ServerPort:         "80",
	DBFile:             "/data/scn.sqlite3",
	RequestTimeout:     16 * time.Second,
	ReturnRawErrors:    true,
	FirebaseProjectID:  "simplecloudnotifier-ea7ef",
	FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
	FirebasePrivKeyID:  "5bfab19fca25034e87c5b3bd1a4334499d2d1f85",
	FirebaseClientMail: "firebase-adminsdk-42grv@simplecloudnotifier-ea7ef.iam.gserviceaccount.com",
	FirebasePrivateKey: "TODO",
}

var configProd = Config{
	Namespace:          "production",
	GinDebug:           false,
	ServerIP:           "0.0.0.0",
	ServerPort:         "80",
	DBFile:             "/data/scn.sqlite3",
	RequestTimeout:     16 * time.Second,
	ReturnRawErrors:    false,
	FirebaseProjectID:  "simplecloudnotifier-ea7ef",
	FirebaseTokenURI:   "https://oauth2.googleapis.com/token",
	FirebasePrivKeyID:  "5bfab19fca25034e87c5b3bd1a4334499d2d1f85",
	FirebaseClientMail: "firebase-adminsdk-42grv@simplecloudnotifier-ea7ef.iam.gserviceaccount.com",
	FirebasePrivateKey: "TODO",
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
