package ginresp

import "github.com/gin-gonic/gin"

type WHandlerFunc func(*gin.Context) HTTPResponse

func Wrap(fn WHandlerFunc) gin.HandlerFunc {

	return func(context *gin.Context) {

		wrap := fn(context)

		if context.Writer.Written() {
			panic("Writing in WrapperFunc is not supported")
		}

		wrap.Write(context)

	}

}
