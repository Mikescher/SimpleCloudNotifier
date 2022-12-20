package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"time"
)

type WHandlerFunc func(*gin.Context) HTTPResponse

func Wrap(fn WHandlerFunc) gin.HandlerFunc {

	maxRetry := scn.Conf.RequestMaxRetry
	retrySleep := scn.Conf.RequestRetrySleep

	return func(g *gin.Context) {

		reqctx := g.Request.Context()

		if g.Request.Body != nil {
			g.Request.Body = dataext.NewBufferedReadCloser(g.Request.Body)
		}

		for ctr := 1; ; ctr++ {

			wrap := fn(g)

			if g.Writer.Written() {
				panic("Writing in WrapperFunc is not supported")
			}

			if ctr < maxRetry && isSqlite3Busy(wrap) {
				log.Warn().Int("counter", ctr).Str("url", g.Request.URL.String()).Msg("Retry request (ErrBusy)")

				err := resetBody(g)
				if err != nil {
					panic(err)
				}

				time.Sleep(retrySleep)
				continue
			}

			if reqctx.Err() == nil {
				wrap.Write(g)
			}

			return
		}

	}

}

func resetBody(g *gin.Context) error {
	if g.Request.Body == nil {
		return nil
	}

	err := g.Request.Body.(dataext.BufferedReadCloser).Reset()
	if err != nil {
		return err
	}

	return nil
}

func isSqlite3Busy(r HTTPResponse) bool {
	if errwrap, ok := r.(*errorHTTPResponse); ok && errwrap != nil {
		if s3err, ok := (errwrap.error).(sqlite3.Error); ok {
			if s3err.Code == sqlite3.ErrBusy {
				return true
			}
		}
	}
	return false
}
