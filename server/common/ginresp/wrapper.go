package ginresp

import "github.com/gin-gonic/gin"

type WHandlerFunc func(*gin.Context) HTTPResponse

func Wrap(fn WHandlerFunc) gin.HandlerFunc {

	return func(g *gin.Context) {

		wrap := fn(g)

		if g.Writer.Written() {
			panic("Writing in WrapperFunc is not supported")
		}

		wrap.Write(g)

	}

}
