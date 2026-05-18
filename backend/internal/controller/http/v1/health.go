package v1

import (
	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	shared.Success(c, gin.H{
		"status":  "ok",
		"version": "0.1.0",
	})
}
