package util

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitTests() {
	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05.000 Z07:00",
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
