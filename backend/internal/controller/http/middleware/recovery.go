package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recovery(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("error", r).
					Str("stack", string(debug.Stack())).
					Str("path", c.Request.URL.Path).
					Msg("Panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    5000,
					"message": "Internal server error",
					"data":    nil,
				})
			}
		}()
		c.Next()
	}
}
