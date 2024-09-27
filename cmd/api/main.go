package main

import (
	"fmt"
	"go-template/internal/auth"
	"go-template/internal/shared"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/internal/shared/interfaces/http"
	"go-template/internal/shared/middleware"
	"log"

	_ "go-template/docs"
	"go-template/internal/config"
	sharedConfig "go-template/internal/shared/config"

	"github.com/gin-gonic/gin"
)

// @title Go Template API Documentation
// @version 1.0
// @description This is a sample server for Go Template API.
// @BasePath /
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Authorization token
func main() {
	loadAppConfig()

	db := initDatabase()
	defer db.Close()

	logger := logger.NewLogrusLogger("./logs")

	authModule := auth.NewModule(db, logger)

	server := http.NewServer()

	server.AddMiddlewares(
		middleware.NewRequestLoggerMiddleware(logger).Handler(),
		middleware.RemovePayloadForMethodNotAllowed(),
		gin.Recovery(),
	)

	server.AddModules(
		shared.NewModule(db),
		authModule,
	)

	serveAndListen(server)
}

func loadAppConfig() {
	if err := sharedConfig.Load(&config.App); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func initDatabase() database.BaseDatabase {
	postgres := database.NewPostgresDatabase(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			config.App.Database.Username,
			config.App.Database.Password,
			config.App.Database.Host,
			config.App.Database.Port,
			config.App.Database.Name,
		),
	)

	return postgres
}

func serveAndListen(server *http.Server) {
	log.Println("Starting server on :" + fmt.Sprint(config.App.Server.Port))
	if err := server.Start(":" + fmt.Sprint(config.App.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
