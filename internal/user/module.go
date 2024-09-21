package user

import (
	"database/sql"
	"go-template/internal/auth"
	"go-template/internal/user/application"
	"go-template/internal/user/domain"
	"go-template/internal/user/infrastructure"
	"go-template/internal/user/interfaces"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler    *interfaces.UserHandler
	authModule *auth.Module
}

func NewModule(db *sql.DB, authModule *auth.Module) *Module {
	userRepo := infrastructure.NewPostgresUserRepository(db)
	userDomainService := domain.NewUserService(userRepo)
	userAppService := application.NewUserApplicationService(userDomainService)
	userHandler := interfaces.NewUserHandler(userAppService)

	return &Module{
		handler:    userHandler,
		authModule: authModule,
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")
	{
		userRouter := apiV1.Group("/users")
		userRouter.Use(m.authModule.GetJWTAuthMiddleware())
		{
			userRouter.GET("/:id", m.handler.GetUser)
			userRouter.PUT("/:id", m.handler.UpdateUser)
		}
	}

}
