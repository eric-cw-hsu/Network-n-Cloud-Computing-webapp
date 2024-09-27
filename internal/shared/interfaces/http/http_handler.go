package http

import (
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/shared/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type SharedHandler struct {
	db     database.BaseDatabase
	logger logger.Logger
}

func NewSharedHandler(db database.BaseDatabase, logger logger.Logger) *SharedHandler {
	return &SharedHandler{
		db:     db,
		logger: logger,
	}
}

// @Summary Database health check
// @Description Check if the database is healthy
// @Tags shared
// @Produce json
// @Success 200
// @Router /healthz [get]
func (h *SharedHandler) Healthz(c *gin.Context) {
	// add no-cache headers
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	// check if there is any payload in the request
	if c.Request.ContentLength > 0 {
		c.Status(400)
		return
	}

	if err := h.db.CheckDBConnection(); err != nil {
		h.logger.Error("Database is not healthy ", err)
		c.Status(503)
		return
	}

	c.Status(200)
	return
}
