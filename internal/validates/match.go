package validates

import "time"

type MatchCreateValidate struct {
	MaleClientID   uint64    `json:"male_client_id" binding:"required"`
	FemaleClientID uint64    `json:"female_client_id" binding:"required"`
	MatchDate      time.Time `json:"match_date"`
	MatchScore     int       `json:"match_score"`
	Remark         string    `json:"remark"`
}

type MatchUpdateValidate struct {
	ID     uint64 `json:"id" binding:"required"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type FollowUpCreateValidate struct {
	MatchRecordID  uint64    `json:"match_record_id" binding:"required"`
	FollowUpDate   time.Time `json:"follow_up_date"`
	Method         string    `json:"method" binding:"required"`
	Content        string    `json:"content" binding:"required"`
	Feedback       string    `json:"feedback"`
	Satisfaction   int8      `json:"satisfaction"`
	Attachments    string    `json:"attachments"`
	NextFollowUpAt time.Time `json:"next_follow_up_at"`
}

type MatchListValidate struct {
	Paginate
	MaleName   string `json:"male_name" form:"male_name"`
	FemaleName string `json:"female_name" form:"female_name"`
	Status     int8   `json:"status" form:"status"`
}
