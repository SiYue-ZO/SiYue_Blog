package app

import (
	"net/http"

	"SiYue_Blog/internal/config"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(cfg *config.Config, engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    cfg.Server.Address,
		Handler: engine,
	}
}
