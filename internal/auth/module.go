package auth

import (
	"database/sql"
	"go-template/internal/auth/application"
	"go-template/internal/auth/config"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/auth/infrastructure"
	"go-template/internal/auth/interfaces/http"
	"go-template/internal/auth/interfaces/http/middleware"
	sharedConfig "go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler    *http.AuthHandler
	jwtService *jwt.JWTService
	authConfig *config.AuthConfig
}

func NewModule(db *sql.DB, logger logger.Logger) *Module {
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

	return &Module{
		handler:    authHandler,
		jwtService: jwtService,
		authConfig: authConfig,
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

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register", m.handler.Register)
		apiV1.POST("/login", m.handler.Login)
	}
}
