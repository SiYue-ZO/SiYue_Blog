package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"siyue-blog/pkg/config"
)

type App struct {
	cfg    *config.Config
	log    zerolog.Logger
	engine *gin.Engine
}

func NewApp(cfg *config.Config, log zerolog.Logger, engine *gin.Engine) *App {
	return &App{
		cfg:    cfg,
		log:    log,
		engine: engine,
	}
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:         a.cfg.HTTP.Addr,
		Handler:      a.engine,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		a.log.Info().Str("addr", srv.Addr).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	a.log.Info().Msg("Server stopped")
	return nil
}
