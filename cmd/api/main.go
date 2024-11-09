package main

import (
	"fmt"
	"go-template/internal/auth"
	"go-template/internal/cloudwatch"
	"go-template/internal/s3"
	"go-template/internal/shared"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/internal/shared/interfaces/http"
	"go-template/internal/shared/middleware"
	"go-template/internal/user"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	setServerMode()

	logger := logger.NewLogrusLogger("./logs")
	cloudWatchModule := cloudwatch.NewModule(logger)
	defer cloudWatchModule.Shutdown()

	db := initDatabase(cloudWatchModule)
	defer db.Close()

	s3Module := s3.NewModule(logger, cloudWatchModule)

	authModule := auth.NewModule(db, logger)
	userModule := user.NewModule(db, logger, authModule.GetBasicService(), s3Module)

	server := http.NewServer()

	server.AddMiddlewares(
		middleware.NewRequestLoggerMiddleware(logger, cloudWatchModule).Handler(),
		middleware.RemovePayloadForMethodNotAllowed(),
		gin.Recovery(),
	)

	server.AddModules(
		shared.NewModule(db, logger),
		authModule,
		userModule,
	)

	setupGracefulShutdown(cloudWatchModule, db)

	serveAndListen(server)
}

func setServerMode() {
	if config.App.Environment == "production" || config.App.Environment == "staging" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func loadAppConfig() {
	if err := sharedConfig.Load(&config.App); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func initDatabase(cloudWatchModule cloudwatch.CloudWatchModule) database.BaseDatabase {
	sslmode := "require"
	if config.App.Environment != "production" && config.App.Environment != "staging" {
		sslmode = "disable"
	}

	postgres := database.NewPostgresDatabase(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=3",
			config.App.Database.Username,
			config.App.Database.Password,
			config.App.Database.Host,
			config.App.Database.Port,
			config.App.Database.Name,
			sslmode,
		),
		cloudWatchModule,
	)

	// comment this block due to the healthz check for assignment
	// Ref: https://northeastern.instructure.com/courses/192927/assignments/2459523
	if err := postgres.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return postgres
}

func serveAndListen(server *http.Server) {
	log.Println("Starting server on :" + fmt.Sprint(config.App.Server.Port))
	if err := server.Start("127.0.0.1:" + fmt.Sprint(config.App.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupGracefulShutdown(
	cloudwatchModule cloudwatch.CloudWatchModule,
	database database.BaseDatabase,
) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Printf("\n--------------------------------\n")
		fmt.Println("Shutting down server...")
		cloudwatchModule.Shutdown()
		database.Close()
		os.Exit(0)
	}()
}
