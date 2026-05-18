package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/pkg/config"
	"siyue-blog/pkg/jwtutil"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			shared.Error(c, http.StatusUnauthorized, shared.ErrCodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			shared.Error(c, http.StatusUnauthorized, shared.ErrCodeUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtutil.ValidateToken(cfg.JWT.AccessSecret, parts[1])
		if err != nil {
			shared.Error(c, http.StatusUnauthorized, shared.ErrCodeInvalidToken, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			shared.Error(c, http.StatusForbidden, shared.ErrCodeForbidden, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
