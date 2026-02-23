package biz_omiai

type AIMatchRepo interface {
	GenerateDailyRecommendations() error
	GetDailyStats() (map[string]interface{}, error)
}
