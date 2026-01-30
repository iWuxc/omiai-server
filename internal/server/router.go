package server

import (
	"net/http"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/client"
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
	ClientController *client.Controller
}

func (r *Router) Register() http.Handler {
	g := r.Group("api")
	{
		//r.user_custom_brand(g.Group("user_custom_brand", middleware.Authorization(r.DB, r.Redis)))
		r.banner(g.Group("banner"))
		r.client(g.Group("client"))
	}
	return r
}

func (r *Router) banner(g *gin.RouterGroup) {
	g.GET("/detail", r.BannerController.Detail) // demo
}

func (r *Router) client(g *gin.RouterGroup) {
	g.POST("/create", r.ClientController.Create)
	g.GET("/list", r.ClientController.List)
	g.GET("/detail/:id", r.ClientController.Detail)
	g.GET("/match/:id", r.ClientController.Match)
}
