package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
)

type rateWindow struct {
	count int
	reset time.Time
}

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	const maxTrackedClients = 10000
	var mu sync.Mutex
	clients := make(map[string]rateWindow)
	lastCleanup := time.Now()
	return func(c *gin.Context) {
		now := time.Now()
		key := c.ClientIP()
		mu.Lock()
		entry, ok := clients[key]
		cleanupInterval := min(window, time.Minute)
		if now.Sub(lastCleanup) >= cleanupInterval {
			for client, item := range clients {
				if !now.Before(item.reset) {
					delete(clients, client)
				}
			}
			lastCleanup = now
		}
		if !ok && len(clients) >= maxTrackedClients {
			mu.Unlock()
			c.Header("Retry-After", "60")
			shared.Error(c, http.StatusTooManyRequests, shared.ErrCodeBadRequest, "too many clients")
			c.Abort()
			return
		}
		if !ok || now.After(entry.reset) {
			entry = rateWindow{reset: now.Add(window)}
		}
		entry.count++
		clients[key] = entry
		remaining := maxRequests - entry.count
		mu.Unlock()
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		if entry.count > maxRequests {
			c.Header("Retry-After", strconv.Itoa(max(1, int(time.Until(entry.reset).Seconds()))))
			shared.Error(c, http.StatusTooManyRequests, shared.ErrCodeBadRequest, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
