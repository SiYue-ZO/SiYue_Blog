package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/pkg/config"
)

type HealthHandler struct {
	db      *gorm.DB
	version string
}

func NewHealthHandler(db *gorm.DB, cfg *config.Config) *HealthHandler {
	return &HealthHandler{db: db, version: cfg.App.Version}
}

func (h *HealthHandler) Check(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		shared.Error(c, http.StatusServiceUnavailable, shared.ErrCodeDatabase, "database unavailable")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		shared.Error(c, http.StatusServiceUnavailable, shared.ErrCodeDatabase, "database unavailable")
		return
	}
	shared.Success(c, gin.H{
		"status":  "ok",
		"version": h.version,
	})
}
