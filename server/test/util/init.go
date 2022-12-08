package util

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func InitTests() {
	log.Logger = createLogger(createConsoleWriter())

	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func createConsoleWriter() *zerolog.ConsoleWriter {
	return &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05.000 Z07:00",
	}
}

func createLogger(cw io.Writer) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	return logger
}
