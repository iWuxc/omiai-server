package server

import (
	"net/http"
	"omiai-server/internal/controller/ai"
	"omiai-server/internal/controller/auth"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/match"
	"omiai-server/internal/controller/reminder"
	"omiai-server/internal/data"
	"omiai-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/redis"
)

// Router .
type Router struct {
	*gin.Engine
	DB                 *data.DB
	Redis              *redis.Redis
	AIController       *ai.Controller
	AuthController     *auth.Controller
	BannerController   *banner.Controller
	ClientController   *client.Controller
	CommonController   *common.Controller
	MatchController    *match.Controller
	ReminderController *reminder.Controller
}

func (r *Router) Register() http.Handler {
	g := r.Group("api")
	{
		r.auth(g.Group("auth"))

		// 需要登录的接口
		authGroup := g.Group("", middleware.Authorization(r.DB, r.Redis))
		{
			r.ai(authGroup.Group("ai"))
			r.banner(authGroup.Group("banner"))
			r.client(authGroup.Group("clients")) // Renamed from "client" to "clients" for V2
			r.common(authGroup.Group("common"))
			r.match(authGroup.Group("couples")) // Renamed from "match" to "couples" for V2
			r.reminder(authGroup.Group("reminders"))
			authGroup.GET("/user/info", r.AuthController.GetUserInfo)
			authGroup.POST("/user/change_password", r.AuthController.ChangePassword)
		}
	}
	// Serve static files for uploads
	r.Static("/uploads", "./runtime/uploads")

	return r
}

func (r *Router) auth(g *gin.RouterGroup) {
	g.POST("/send_sms", r.AuthController.SendSms)
	g.POST("/login/h5", r.AuthController.H5Login)
	g.POST("/login/wx", r.AuthController.WxLogin)
}

func (r *Router) ai(g *gin.RouterGroup) {
	g.POST("/analyze", r.AIController.AnalyzeMatch)
	g.POST("/ice-breaker", r.AIController.GetIceBreaker)
}

func (r *Router) common(g *gin.RouterGroup) {
	g.POST("/upload", r.CommonController.Upload)
}

func (r *Router) match(g *gin.RouterGroup) {
	g.GET("/list", r.MatchController.List)
	g.POST("/create", r.MatchController.Create)
	g.POST("/confirm", r.MatchController.Confirm)   // V2: Direct Confirm
	g.POST("/dissolve", r.MatchController.Dissolve) // V2: Dissolve Match
	g.POST("/update_status", r.MatchController.UpdateStatus)
	g.GET("/followup/list", r.MatchController.ListFollowUps)
	g.POST("/followup/create", r.MatchController.CreateFollowUp)
	g.GET("/reminders", r.MatchController.GetReminders)
	g.GET("/status/history", r.MatchController.GetStatusHistory)
	g.GET("/stats", r.MatchController.Stats)
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
	g.POST("/update", r.ClientController.Update)
	g.DELETE("/delete/:id", r.ClientController.Delete)
	g.GET("/list", r.ClientController.List)
	g.GET("/detail/:id", r.ClientController.Detail)
	g.GET("/match/:id", r.ClientController.MatchV2) // Upgrade to V2
	// V2: New Candidates & Compare Interfaces
	g.GET("/:id/candidates", r.MatchController.GetCandidates)
	g.GET("/:id/compare/:candidateId", r.MatchController.Compare)

	// Phase 1: Claim/Release (Hidden for Single Mode but kept for compatibility)
	g.POST("/claim", r.ClientController.Claim)
	g.POST("/release", r.ClientController.Release)

	// Import
	g.POST("/import/analyze", r.ClientController.ImportAnalyze)
	g.POST("/import/batch", r.ClientController.ImportBatch)
}

func (r *Router) reminder(g *gin.RouterGroup) {
	g.GET("/list", r.ReminderController.List)
	g.GET("/today", r.ReminderController.TodayList)
	g.GET("/pending", r.ReminderController.PendingList)
	g.GET("/stats", r.ReminderController.Stats)
	g.POST("/read", r.ReminderController.MarkAsRead)
	g.POST("/done", r.ReminderController.MarkAsDone)
	g.DELETE("/delete", r.ReminderController.Delete)
}
