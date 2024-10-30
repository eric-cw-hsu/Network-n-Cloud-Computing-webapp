package user

import (
	"go-template/internal/auth/domain/basic"
	"go-template/internal/auth/interfaces/http/middleware"
	"go-template/internal/s3"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/internal/user/application"
	"go-template/internal/user/domain"
	"go-template/internal/user/infrastructure"
	"go-template/internal/user/interfaces/http"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *http.UserHandler
	basicService *basic.BasicService
	s3Module     s3.S3Module
}

func NewModule(db database.BaseDatabase, logger logger.Logger, basicService *basic.BasicService, s3Module s3.S3Module) *Module {
	userRepository := infrastructure.NewPostgresUserRepository(db)
	userService := domain.NewUserService(s3Module)

	userApplicationService := application.NewUserApplicationService(logger, userService, userRepository)
	userHandler := http.NewUserHandler(userApplicationService, s3Module)

	return &Module{
		handler:      userHandler,
		basicService: basicService,
		s3Module:     s3Module,
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine) {

	userRouter := router.Group("/v1/user")
	userRouter.Use(middleware.BasicAuthMiddleware(m.basicService))
	{
		userRouter.POST("/self/pic", m.handler.UploadProfilePic)
		userRouter.GET("/self/pic", m.handler.GetProfilePic)
		userRouter.DELETE("/self/pic", m.handler.DeleteProfilePic)
	}
}
