package server

import (
	"net/http"
	"omiai-server/internal/controller/ai"
	"omiai-server/internal/controller/auth"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/c_auth"
	"omiai-server/internal/controller/c_client"
	"omiai-server/internal/controller/c_event"
	"omiai-server/internal/controller/c_interact"
	"omiai-server/internal/controller/c_pay"
	"omiai-server/internal/controller/c_recommend"
	"omiai-server/internal/controller/china_region"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/dashboard"
	"omiai-server/internal/controller/event"
	"omiai-server/internal/controller/match"
	"omiai-server/internal/controller/reminder"
	"omiai-server/internal/controller/template"
	"omiai-server/internal/controller/tenant"
	"omiai-server/internal/data"
	"omiai-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/redis"
)

// Router .
type Router struct {
	*gin.Engine
	DB                    *data.DB
	Redis                 *redis.Redis
	AIController          *ai.Controller
	AuthController        *auth.Controller
	CAuthController       *c_auth.Controller
	CClientController     *c_client.Controller
	CRecommendController  *c_recommend.Controller
	CInteractController   *c_interact.Controller
	CPayController        *c_pay.Controller
	CEventController      *c_event.Controller
	BannerController      *banner.Controller
	ChinaRegionController *china_region.Controller
	ClientController      *client.Controller
	CommonController      *common.Controller
	TemplateController    *template.Controller
	ReminderController    *reminder.Controller
	DashboardController   *dashboard.Controller
	MatchController       *match.Controller
	TenantController      *tenant.Controller
	EventController       *event.Controller
}

func (r *Router) Register() http.Handler {
	g := r.Group("api")
	{
		r.auth(g.Group("auth"))
		// 不需要登录的接口
		g.GET("/china_region/provinces", r.ChinaRegionController.GetProvinces)
		g.GET("/china_region/cities", r.ChinaRegionController.GetCities)
		g.GET("/china_region/districts", r.ChinaRegionController.GetDistricts)
		g.GET("/china_region/hot", r.ChinaRegionController.GetHotCities)
		g.GET("/china_region/search", r.ChinaRegionController.Search)

		// C端接口 (小程序)
		cGroup := g.Group("c")
		{
			// C端鉴权与回调
			cGroup.POST("/auth/wx_login", r.CAuthController.WxLogin)
			cGroup.POST("/pay/wechat_notify", r.CPayController.WechatNotify)

			// C端需要登录的接口
			cAuthGroup := cGroup.Group("", middleware.CClientAuthorization(r.DB, r.Redis))
			{
				cAuthGroup.GET("/profile/mine", r.CClientController.GetMine)
				cAuthGroup.POST("/profile/update", r.CClientController.UpdateMine)
				cAuthGroup.POST("/profile/verify_realname", r.CClientController.VerifyRealName)

				// 推荐流
				cAuthGroup.GET("/recommend/daily", r.CRecommendController.DailyRecommend)
				cAuthGroup.GET("/recommend/detail/:id", r.CRecommendController.Detail)

				// 互动
				cAuthGroup.POST("/interact/like", r.CInteractController.Like)
				cAuthGroup.GET("/interact/received_likes", r.CInteractController.GetReceivedLikes)
				cAuthGroup.GET("/interact/mutual_matches", r.CInteractController.GetMutualMatches)

				// 商业化与支付
				cAuthGroup.POST("/pay/recharge", r.CPayController.Recharge)
				cAuthGroup.POST("/pay/unlock", r.CPayController.UnlockProfile)
				cAuthGroup.POST("/pay/vip", r.CPayController.BuyVip)

				// 线下活动
				cAuthGroup.GET("/event/list", r.CEventController.List)
				cAuthGroup.POST("/event/register", r.CEventController.Register)
			}
		}

		// 需要登录的接口
		authGroup := g.Group("", middleware.Authorization(r.DB, r.Redis))
		{
			// SaaS 超管接口 (只允许 admin 访问)
			saasGroup := authGroup.Group("saas")
			{
				saasGroup.POST("/tenant", r.TenantController.Create)
				saasGroup.GET("/tenant/list", r.TenantController.List)
			}

			r.client(authGroup.Group("client"))
			r.banner(authGroup.Group("banner"))
			r.client(authGroup.Group("clients")) // Renamed from "client" to "clients" for V2
			r.common(authGroup.Group("common"))
			r.dashboard(authGroup.Group("dashboard"))
			r.match(authGroup.Group("couples")) // Renamed from "match" to "couples" for V2
			r.reminder(authGroup.Group("reminders"))
			r.template(authGroup.Group("templates"))

			// B端活动管理
			eventGroup := authGroup.Group("events")
			{
				eventGroup.POST("/publish", r.EventController.Publish)
				eventGroup.GET("/list", r.EventController.List)
			}

			// 认证相关接口（需要登录）
			authGroup.GET("/auth/codes", r.AuthController.GetAccessCodes)
			authGroup.GET("/user/info", r.AuthController.GetUserInfo)
			authGroup.POST("/user/change_password", r.AuthController.ChangePassword)
		}
	}
	// Serve static files for uploads
	r.Static("/uploads", "./runtime/uploads")
	// Serve H5 frontend
	r.Static("/h5", "./web")

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

func (r *Router) dashboard(g *gin.RouterGroup) {
	g.GET("/stats", r.DashboardController.Stats)
	g.GET("/todos", r.DashboardController.GetTodos)
	g.GET("/leads", r.DashboardController.GetInteractionLeads)
	g.POST("/leads/status", r.DashboardController.UpdateInteractionStatus)
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
	g.POST("/verify/:id", r.ClientController.Verify)
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
	g.GET("/list", r.ReminderController.ListPendingTasks) // Default to pending tasks
	// g.GET("/today", r.ReminderController.TodayList) // Removed
	// g.GET("/pending", r.ReminderController.PendingList) // Removed
	// g.GET("/stats", r.ReminderController.Stats) // Removed
	// g.POST("/read", r.ReminderController.MarkAsRead) // Removed
	g.POST("/done/:id", r.ReminderController.CompleteTask) // Updated
	// g.DELETE("/delete", r.ReminderController.Delete) // Removed

	// New routes
	g.GET("/rules", r.ReminderController.ListRules)
	g.POST("/rules", r.ReminderController.CreateRule)
	g.GET("/tasks/pending", r.ReminderController.ListPendingTasks)
	g.POST("/tasks/:id/complete", r.ReminderController.CompleteTask)
	g.POST("/generate", r.ReminderController.CheckAndGenerateTasks)
}

func (r *Router) template(g *gin.RouterGroup) {
	g.POST("", r.TemplateController.Create)
	g.GET("", r.TemplateController.List)
	g.PUT("/:id", r.TemplateController.Update)
	g.DELETE("/:id", r.TemplateController.Delete)
	g.POST("/:id/use", r.TemplateController.Use)
}
