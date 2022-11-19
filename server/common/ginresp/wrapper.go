package ginresp

import "github.com/gin-gonic/gin"

type WHandlerFunc func(*gin.Context) HTTPResponse

func Wrap(fn WHandlerFunc) gin.HandlerFunc {

	return func(g *gin.Context) {

		reqctx := g.Request.Context()

		wrap := fn(g)

		if g.Writer.Written() {
			panic("Writing in WrapperFunc is not supported")
		}

		if reqctx.Err() == nil {
			wrap.Write(g)
		}

	}

}
