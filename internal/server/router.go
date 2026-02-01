package server

import (
	"net/http"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/match"
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
	CommonController *common.Controller
	MatchController  *match.Controller
}

func (r *Router) Register() http.Handler {
	g := r.Group("api")
	{
		//r.user_custom_brand(g.Group("user_custom_brand", middleware.Authorization(r.DB, r.Redis)))
		r.banner(g.Group("banner"))
		r.client(g.Group("client"))
		r.common(g.Group("common"))
		r.match(g.Group("match"))
	}
	// Serve static files for uploads
	r.Static("/uploads", "./runtime/uploads")

	return r
}

func (r *Router) common(g *gin.RouterGroup) {
	g.POST("/upload", r.CommonController.Upload)
}

func (r *Router) match(g *gin.RouterGroup) {
	g.GET("/list", r.MatchController.List)
	g.POST("/create", r.MatchController.Create)
	g.POST("/update_status", r.MatchController.UpdateStatus)
	g.GET("/followup/list", r.MatchController.ListFollowUps)
	g.POST("/followup/create", r.MatchController.CreateFollowUp)
	g.GET("/reminders", r.MatchController.GetReminders)
}

func (r *Router) banner(g *gin.RouterGroup) {
	g.GET("/list", r.BannerController.List)
	g.GET("/detail", r.BannerController.Detail) // demo
	g.POST("/create", r.BannerController.Create)
	g.POST("/update", r.BannerController.Update)
	g.DELETE("/delete/:id", r.BannerController.Delete)
}

func (r *Router) client(g *gin.RouterGroup) {
	g.GET("/stats", r.ClientController.Stats)
	g.POST("/create", r.ClientController.Create)
	g.GET("/list", r.ClientController.List)
	g.GET("/detail/:id", r.ClientController.Detail)
	g.GET("/match/:id", r.ClientController.Match)
}
