package di

import (
	"SiYue_Blog/internal/app"
	"SiYue_Blog/internal/config"
	"SiYue_Blog/internal/db"
	"SiYue_Blog/internal/handler"
	"SiYue_Blog/internal/repository"
	"SiYue_Blog/internal/router"
	"SiYue_Blog/internal/service"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	config.New,
	db.NewMySQL,
	repository.NewHealthRepository,
	service.NewHealthService,
	handler.NewHealthHandler,
	router.NewEngine,
	app.NewHTTPServer,
)
