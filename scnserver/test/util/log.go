package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var buflogger *BufferWriter = nil

func SetBufLogger() {
	buflogger = &BufferWriter{cw: createConsoleWriter()}
	log.Logger = createLogger(buflogger)
	gin.SetMode(gin.ReleaseMode)
}

func ClearBufLogger(dump bool) {
	size := len(buflogger.buffer)
	if dump {
		buflogger.Dump()
	}
	log.Logger = createLogger(createConsoleWriter())
	buflogger = nil
	gin.SetMode(gin.TestMode)
	if !dump {
		log.Info().Msgf("Suppressed %d logmessages / printf-statements", size)
	}
}

func TPrintf(lvl zerolog.Level, format string, a ...any) {
	if zerolog.GlobalLevel() <= lvl {
		if buflogger != nil {
			buflogger.Printf(format, a...)
		} else {
			fmt.Printf(format, a...)
		}
	}
}

func TPrintln(lvl zerolog.Level, a ...any) {
	if zerolog.GlobalLevel() <= lvl {
		if buflogger != nil {
			buflogger.Println(a...)
		} else {
			fmt.Println(a...)
		}
	}
}
