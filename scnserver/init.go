package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var callerRoot = ""

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	callerRoot = path.Dir(file)
}

func Init(cfg Config) {
	cw := zerolog.ConsoleWriter{
		Out:          os.Stdout,
		TimeFormat:   "2006-01-02 15:04:05 Z07:00",
		FormatCaller: formatCaller,
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	if cfg.GinDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	zerolog.SetGlobalLevel(cfg.LogLevel)

	log.Debug().Msg("Initialized")
}

func formatCaller(i any) string {
	const (
		colorBlack = iota + 30
		colorRed
		colorGreen
		colorYellow
		colorBlue
		colorMagenta
		colorCyan
		colorWhite

		colorBold     = 1
		colorDarkGray = 90
	)

	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		if rel, err := filepath.Rel(callerRoot, c); err == nil {
			c = rel
		}
		c = colorize(c, colorBold, false) + colorize(" >", colorCyan, false)
	}
	return c
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
