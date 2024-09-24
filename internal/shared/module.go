package shared

import (
	"go-template/internal/shared/interfaces/http"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *http.SharedHandler
}

func NewModule() *Module {
	handler := http.NewSharedHandler()

	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	router.GET("/healthz", m.handler.Healthz)
}
