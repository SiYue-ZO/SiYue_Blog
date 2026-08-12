package http

import (
	"net/http"
	"time"

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
	contactHandler *v1.ContactHandler,
	adminHandler *v1.AdminHandler,
	settingsHandler *v1.SettingsHandler,
	homePageHandler *v1.HomePageHandler,
) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.LimitBody(cfg.HTTP.MaxBodyBytes))
	r.Use(middleware.CORS(cfg))

	// Public API
	api := r.Group("/api/v1")
	{
		// Health
		api.GET("/health", healthHandler.Check)
		api.GET("/settings", settingsHandler.Get)
		api.GET("/home", homePageHandler.GetPublished)

		// Auth
		api.POST("/auth/login", middleware.RequireJSON(), middleware.RateLimit(10, 15*time.Minute), authHandler.Login)
		if cfg.Auth.AllowRegistration {
			api.POST("/auth/register", middleware.RequireJSON(), middleware.RateLimit(5, time.Hour), authHandler.Register)
		}
		api.POST("/auth/refresh", middleware.RateLimit(30, time.Hour), authHandler.RefreshToken)
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", middleware.JWTAuth(cfg), authHandler.Me)

		// Posts
		api.GET("/posts", postHandler.List)
		api.GET("/posts/featured", postHandler.Featured)
		api.GET("/posts/:slug", postHandler.GetBySlug)

		// Tags
		api.GET("/tags", tagHandler.List)

		// Projects
		api.GET("/projects", projectHandler.List)
		api.GET("/projects/featured", projectHandler.Featured)

		// Contact
		api.POST("/contact", middleware.RequireJSON(), middleware.RateLimit(5, time.Hour), contactHandler.Create)
	}

	// Admin API (requires JWT + admin role)
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(cfg))
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/stats", adminHandler.Stats)
		admin.GET("/settings", settingsHandler.Get)
		admin.PATCH("/settings", middleware.RequireJSON(), settingsHandler.Update)
		admin.GET("/home", homePageHandler.GetAdmin)
		admin.PUT("/home/draft", middleware.RequireJSON(), homePageHandler.SaveDraft)
		admin.POST("/home/publish", homePageHandler.Publish)
		admin.POST("/home/reset", homePageHandler.ResetDraft)
		admin.POST("/password", middleware.RequireJSON(), authHandler.ChangePassword)
		admin.GET("/posts", postHandler.AdminList)
		admin.GET("/posts/:id", postHandler.AdminGet)
		admin.POST("/posts", middleware.RequireJSON(), postHandler.Create)
		admin.PUT("/posts/:id", middleware.RequireJSON(), postHandler.Update)
		admin.DELETE("/posts/:id", postHandler.Delete)

		admin.GET("/tags", tagHandler.List)
		admin.POST("/tags", middleware.RequireJSON(), tagHandler.Create)
		admin.DELETE("/tags/:id", tagHandler.Delete)

		admin.GET("/projects", projectHandler.AdminList)
		admin.POST("/projects", middleware.RequireJSON(), projectHandler.Create)
		admin.PUT("/projects/:id", middleware.RequireJSON(), projectHandler.Update)
		admin.DELETE("/projects/:id", projectHandler.Delete)

		admin.GET("/messages", contactHandler.List)
		admin.PATCH("/messages/:id", middleware.RequireJSON(), contactHandler.UpdateStatus)
		admin.DELETE("/messages/:id", contactHandler.Delete)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 4000, "message": "not found", "data": nil})
	})

	return r
}
