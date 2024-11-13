package auth

import (
	"go-template/internal/auth/application"
	"go-template/internal/auth/config"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/basic"
	"go-template/internal/auth/infrastructure"
	"go-template/internal/auth/interfaces/http"
	"go-template/internal/auth/interfaces/http/middleware"
	"go-template/internal/aws/sns"
	sharedConfig "go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *http.AuthHandler
	basicService *basic.BasicService
	authConfig   *config.AuthConfig
}

func NewModule(db database.BaseDatabase, logger logger.Logger, snsModule sns.SNSModule) *Module {
	// load auth config with viper
	authConfig := loadConfig()

	authRepo := infrastructure.NewPostgresAuthRepository(db)
	authDomainService := domain.NewAuthService(authRepo, logger, authConfig, snsModule)
	authAppService := application.NewAuthApplicationService(authDomainService, logger)
	authHandler := http.NewAuthHandler(authAppService)
	basicService := basic.NewBasicService(authRepo)

	return &Module{
		handler:      authHandler,
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

func (m *Module) GetBasicService() *basic.BasicService {
	return m.basicService
}

func (m *Module) RegisterRoutes(router *gin.Engine) {

	router.GET("/verify", m.handler.VerifyAccount)

	v1User := router.Group("/v1/user")
	{
		v1User.POST("", m.handler.Register)

		// the route below protected by basic auth middleware
		authenticated := v1User.Group("")
		authenticated.Use(middleware.BasicAuthMiddleware(m.basicService))
		{
			authenticated.GET("/resend-verification-email", m.handler.ResendVerification)

			authenticated.Use(middleware.AccountVerificationMiddleware())
			{
				authenticated.GET("/self", m.handler.GetUser)
				authenticated.PUT("/self", m.handler.UpdateUser)
			}
		}
	}
}
