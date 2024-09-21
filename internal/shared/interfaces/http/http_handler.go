package http

import (
	"fmt"
	"go-template/internal/config"
	"go-template/internal/shared/infrastructure/postgres"

	"github.com/gin-gonic/gin"
)

type SharedHandler struct {
}

func NewSharedHandler() *SharedHandler {
	return &SharedHandler{}
}

// @Summary Database health check
// @Description Check if the database is healthy
// @Tags shared
// @Produce json
// @Success 200
// @Router /healthz [get]
func (h *SharedHandler) Healthz(c *gin.Context) {
	// TODO: reduce the dependency on postgres package (should implement a database interface)

	// check if there is any payload in the request
	if c.Request.ContentLength > 0 {
		c.Status(400)
		return
	}

	// add no-cache headers
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	dbSourceString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.App.Database.Username,
		config.App.Database.Password,
		config.App.Database.Host,
		config.App.Database.Port,
		config.App.Database.Name,
	)
	if postgres.CheckDBConnection(dbSourceString) != nil {
		c.Status(503)
		return
	}

	c.Status(200)
	return
}
