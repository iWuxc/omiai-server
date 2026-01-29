package server

import (
	"net/http"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/redis"
)

// Router .
type Router struct {
	*gin.Engine
	DB               *data.DB
	Redis            *redis.Redis
	BannerController *banner.Controller
}

func (r *Router) Register() http.Handler {
	g := r.Group("api")
	{
		//r.user_custom_brand(g.Group("user_custom_brand", middleware.Authorization(r.DB, r.Redis)))
		r.banner(g.Group("banner"))
	}
	return r
}

func (r *Router) banner(g *gin.RouterGroup) {
	g.GET("/detail", r.BannerController.Detail) // demo
}
