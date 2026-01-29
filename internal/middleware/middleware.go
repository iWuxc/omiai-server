package middleware

import (
	"omiai-server/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderMiddlewareSet = wire.NewSet(
	NewMiddleWares,
)

// NewMiddleWares Global Middlewares.
func NewMiddleWares(db *data.DB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		metrics(),
		SetPolicy(),
	}
}
