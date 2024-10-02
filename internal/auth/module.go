package auth

import (
	"go-template/internal/auth/application"
	"go-template/internal/auth/config"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/basic"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/auth/infrastructure"
	"go-template/internal/auth/interfaces/http"
	"go-template/internal/auth/interfaces/http/middleware"
	sharedConfig "go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *http.AuthHandler
	jwtService   *jwt.JWTService
	basicService *basic.BasicService
	authConfig   *config.AuthConfig
}

func NewModule(db database.BaseDatabase, logger logger.Logger) *Module {
	// load auth config with viper
	authConfig := loadConfig()

	jwtConfig := &jwt.JWTConfig{
		JWTSecret:       authConfig.Auth.JWTSecret,
		TokenExpiration: authConfig.Auth.TokenExpiration,
	}
	jwtService := jwt.NewJWTService(jwtConfig)

	authRepo := infrastructure.NewPostgresAuthRepository(db)
	authDomainService := domain.NewAuthService(authRepo, logger)
	authAppService := application.NewAuthApplicationService(authDomainService, jwtService, logger)
	authHandler := http.NewAuthHandler(authAppService)
	basicService := basic.NewBasicService(authRepo)

	return &Module{
		handler:      authHandler,
		jwtService:   jwtService,
		authConfig:   authConfig,
		basicService: basicService,
	}
}

func loadConfig() *config.AuthConfig {
	// load auth config with viper
	var authConfig *config.AuthConfig
	if err := sharedConfig.Load(&authConfig); err != nil {
		log.Fatalf("Failed to load auth config: %v", err)
	}

	return authConfig
}

func (m *Module) GetJWTAuthMiddleware() gin.HandlerFunc {
	return middleware.JWTAuthMiddleware(m.jwtService)
}

func (m *Module) RegisterRoutes(router *gin.Engine) {

	V1 := router.Group("/v1")
	{
		V1.POST("/user", m.handler.Register)
		// apiV1.POST("/login", m.handler.Login)

		// the route below protected by basic auth middleware
		authenticated := V1.Group("/")
		authenticated.Use(middleware.BasicAuthMiddleware(m.basicService))
		{
			authenticated.GET("/user", m.handler.GetUser)
		}
	}
}
