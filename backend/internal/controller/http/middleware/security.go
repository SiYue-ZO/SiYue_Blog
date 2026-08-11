package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if len(id) < 8 || len(id) > 128 {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err == nil {
				id = hex.EncodeToString(bytes)
			}
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}

func LimitBody(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil && c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

func RequireJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Type") != "" && c.ContentType() != "application/json" {
			shared.Error(c, http.StatusUnsupportedMediaType, shared.ErrCodeBadRequest, "Content-Type must be application/json")
			c.Abort()
			return
		}
		c.Next()
	}
}
