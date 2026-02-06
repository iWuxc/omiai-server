package validates

// AIAnalyzeValidate AI分析请求参数
type AIAnalyzeValidate struct {
	ClientAID uint64 `json:"client_a_id" binding:"required"` // 客户A ID
	ClientBID uint64 `json:"client_b_id" binding:"required"` // 客户B ID
}
