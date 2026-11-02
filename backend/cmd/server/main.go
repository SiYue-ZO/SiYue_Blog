package main

import (
	"context"
	"log"
	"time"

	"siyue-blog/internal/app"
	"siyue-blog/internal/bootstrap"
	"siyue-blog/internal/controller/http"
	v1 "siyue-blog/internal/controller/http/v1"
	"siyue-blog/internal/repo/persistence"
	authUC "siyue-blog/internal/usecase/auth"
	contactUC "siyue-blog/internal/usecase/contact"
	homePageUC "siyue-blog/internal/usecase/homepage"
	postUC "siyue-blog/internal/usecase/post"
	projectUC "siyue-blog/internal/usecase/project"
	settingsUC "siyue-blog/internal/usecase/settings"
	tagUC "siyue-blog/internal/usecase/tag"
	"siyue-blog/migrations"
	"siyue-blog/pkg/config"
	pkglogger "siyue-blog/pkg/logger"
	"siyue-blog/pkg/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := pkglogger.NewLogger(cfg.Log.Level, cfg.Log.Format)

	// Initialize database
	db, err := postgres.NewPostgres(cfg.Postgres, cfg.App.Env)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to access database pool")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close database pool")
		}
	}()
	migrationCtx, migrationCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	err = migrations.Up(migrationCtx, db)
	migrationCancel()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Initialize repositories
	userRepo := persistence.NewUserRepo(db)
	sessionRepo := persistence.NewSessionRepo(db)
	postRepo := persistence.NewPostRepo(db)
	tagRepo := persistence.NewTagRepo(db)
	projectRepo := persistence.NewProjectRepo(db)
	contactRepo := persistence.NewContactRepo(db)
	settingsRepo := persistence.NewSettingsRepo(db)
	homePageRepo := persistence.NewHomePageRepo(db)
	bootstrapCtx, bootstrapCancel := context.WithTimeout(context.Background(), 30*time.Second)
	err = bootstrap.EnsureAdmin(bootstrapCtx, userRepo, cfg.Admin)
	bootstrapCancel()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to bootstrap administrator")
	}

	// Initialize use cases
	authUseCase := authUC.NewUseCase(userRepo, sessionRepo, cfg)
	postUseCase := postUC.NewUseCase(postRepo, tagRepo)
	tagUseCase := tagUC.NewUseCase(tagRepo)
	projectUseCase := projectUC.NewUseCase(projectRepo)
	contactUseCase := contactUC.NewUseCase(contactRepo)
	settingsUseCase := settingsUC.NewUseCase(settingsRepo)
	homePageUseCase := homePageUC.NewUseCase(homePageRepo)

	// Initialize handlers
	healthHandler := v1.NewHealthHandler(db, cfg)
	postHandler := v1.NewPostHandler(postUseCase, logger)
	tagHandler := v1.NewTagHandler(tagUseCase)
	projectHandler := v1.NewProjectHandler(projectUseCase)
	authHandler := v1.NewAuthHandler(authUseCase, cfg)
	contactHandler := v1.NewContactHandler(contactUseCase, cfg.JWT.RefreshSecret)
	adminHandler := v1.NewAdminHandler(db)
	settingsHandler := v1.NewSettingsHandler(settingsUseCase)
	homePageHandler := v1.NewHomePageHandler(homePageUseCase)

	// Initialize router
	router := http.NewRouter(cfg, logger, healthHandler, postHandler, tagHandler, projectHandler, authHandler, contactHandler, adminHandler, settingsHandler, homePageHandler)

	// Create and run application
	application := app.NewApp(cfg, logger, router)
	if err := application.Run(); err != nil {
		logger.Fatal().Err(err).Msg("Application failed")
	}
}
