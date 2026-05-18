package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"siyue-blog/internal/controller/http/middleware"
	v1 "siyue-blog/internal/controller/http/v1"
	"siyue-blog/pkg/config"
)

func NewRouter(
	cfg *config.Config,
	log zerolog.Logger,
	healthHandler *v1.HealthHandler,
	postHandler *v1.PostHandler,
	tagHandler *v1.TagHandler,
	projectHandler *v1.ProjectHandler,
	authHandler *v1.AuthHandler,
) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(cfg))

	// Public API
	api := r.Group("/api/v1")
	{
		// Health
		api.GET("/health", healthHandler.Check)

		// Auth
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/refresh", authHandler.RefreshToken)

		// Posts
		api.GET("/posts", postHandler.List)
		api.GET("/posts/featured", postHandler.Featured)
		api.GET("/posts/:slug", postHandler.GetBySlug)

		// Tags
		api.GET("/tags", tagHandler.List)

		// Projects
		api.GET("/projects", projectHandler.List)
		api.GET("/projects/featured", projectHandler.Featured)
	}

	// Admin API (requires JWT + admin role)
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(cfg))
	admin.Use(middleware.AdminOnly())
	{
		// Posts management
		// admin.GET("/posts", adminPostHandler.List)
		// admin.POST("/posts", adminPostHandler.Create)
		// admin.PUT("/posts/:id", adminPostHandler.Update)
		// admin.DELETE("/posts/:id", adminPostHandler.Delete)

		// Tags management
		// admin.GET("/tags", adminTagHandler.List)
		// admin.POST("/tags", adminTagHandler.Create)
		// admin.DELETE("/tags/:id", adminTagHandler.Delete)

		// Projects management
		// admin.GET("/projects", adminProjectHandler.List)
		// admin.POST("/projects", adminProjectHandler.Create)
		// admin.PUT("/projects/:id", adminProjectHandler.Update)
		// admin.DELETE("/projects/:id", adminProjectHandler.Delete)
	}

	return r
}
