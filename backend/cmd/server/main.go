package main

import (
	"log"

	"siyue-blog/internal/app"
	"siyue-blog/internal/controller/http"
	v1 "siyue-blog/internal/controller/http/v1"
	"siyue-blog/internal/repo/persistence"
	authUC "siyue-blog/internal/usecase/auth"
	postUC "siyue-blog/internal/usecase/post"
	projectUC "siyue-blog/internal/usecase/project"
	tagUC "siyue-blog/internal/usecase/tag"
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
	db, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := persistence.NewUserRepo(db)
	postRepo := persistence.NewPostRepo(db)
	tagRepo := persistence.NewTagRepo(db)
	projectRepo := persistence.NewProjectRepo(db)

	// Initialize use cases
	authUseCase := authUC.NewUseCase(userRepo, cfg)
	postUseCase := postUC.NewUseCase(postRepo, tagRepo)
	tagUseCase := tagUC.NewUseCase(tagRepo)
	projectUseCase := projectUC.NewUseCase(projectRepo)

	// Initialize handlers
	healthHandler := v1.NewHealthHandler()
	postHandler := v1.NewPostHandler(postUseCase)
	tagHandler := v1.NewTagHandler(tagUseCase)
	projectHandler := v1.NewProjectHandler(projectUseCase)
	authHandler := v1.NewAuthHandler(authUseCase)

	// Initialize router
	router := http.NewRouter(cfg, logger, healthHandler, postHandler, tagHandler, projectHandler, authHandler)

	// Create and run application
	application := app.NewApp(cfg, logger, router)
	if err := application.Run(); err != nil {
		logger.Fatal().Err(err).Msg("Application failed")
	}
}
