package validates

// ReminderListValidate 提醒列表请求参数
type ReminderListValidate struct {
	Page     int  `form:"page" json:"page"`           // 页码
	PageSize int  `form:"page_size" json:"page_size"` // 每页数量
	IsDone   int8 `form:"is_done" json:"is_done"`     // 是否已完成：-1=全部，0=未完成，1=已完成
	Type     int8 `form:"type" json:"type"`           // 提醒类型
	Priority int8 `form:"priority" json:"priority"`   // 优先级
}

// ReminderIDValidate 提醒ID请求参数
type ReminderIDValidate struct {
	ID uint64 `json:"id" binding:"required"` // 提醒ID
}
