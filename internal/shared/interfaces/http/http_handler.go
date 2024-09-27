package http

import (
	"go-template/internal/shared/infrastructure/database"

	"github.com/gin-gonic/gin"
)

type SharedHandler struct {
	db database.BaseDatabase
}

func NewSharedHandler(db database.BaseDatabase) *SharedHandler {
	return &SharedHandler{
		db: db,
	}
}

// @Summary Database health check
// @Description Check if the database is healthy
// @Tags shared
// @Produce json
// @Success 200
// @Router /healthz [get]
func (h *SharedHandler) Healthz(c *gin.Context) {
	// check if there is any payload in the request
	if c.Request.ContentLength > 0 {
		c.Status(400)
		return
	}

	// add no-cache headers
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	if h.db.CheckDBConnection() != nil {
		c.Status(503)
		return
	}

	c.Status(200)
	return
}
