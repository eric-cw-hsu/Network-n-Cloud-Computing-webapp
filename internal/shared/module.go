package shared

import (
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/interfaces/http"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *http.SharedHandler
}

func NewModule(db database.BaseDatabase) *Module {
	handler := http.NewSharedHandler(db)

	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	router.GET("/healthz", m.handler.Healthz)
}
