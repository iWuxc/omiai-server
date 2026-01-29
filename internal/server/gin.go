package server

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewGinEngine init gin engine .
func NewGinEngine(middlewares []gin.HandlerFunc) *gin.Engine {
	g := gin.New()
	mode := gin.ReleaseMode
	if conf.GetConfig().Debug {
		mode = gin.DebugMode
		g.Use(gin.Logger())
	}
	g.Use(gin.CustomRecovery(middleware.RecoveryFuncCustomer()))
	//g.Use(gin.Recovery())
	g.Use(middlewares...)
	gin.SetMode(mode)
	return g
}
