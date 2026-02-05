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
	ID       uint64 `json:"id" binding:"required"`
	Status   int8   `json:"status" binding:"required"`
	Remark   string `json:"remark"`
	Reason   string `json:"reason"`
	Operator string `json:"operator"` // Optional, or taken from context
}

type FollowUpCreateValidate struct {
	MatchRecordID  uint64 `json:"match_record_id" binding:"required"`
	FollowUpDate   string `json:"follow_up_date"`
	Method         string `json:"method" binding:"required"`
	Content        string `json:"content" binding:"required"`
	Feedback       string `json:"feedback"`
	Satisfaction   int8   `json:"satisfaction"`
	Attachments    string `json:"attachments"`
	NextFollowUpAt string `json:"next_follow_up_at"`
}

type MatchListValidate struct {
	Paginate
	MaleName   string `json:"male_name" form:"male_name"`
	FemaleName string `json:"female_name" form:"female_name"`
	Status     int8   `json:"status" form:"status"`
}

// V2: New Validations

type GetCandidatesValidate struct {
	ClientID uint64 `uri:"id" binding:"required"`
}

type CompareValidate struct {
	ClientID    uint64 `uri:"id" binding:"required"`
	CandidateID uint64 `uri:"candidateId" binding:"required"`
}

type ConfirmMatchValidate struct {
	ClientID    uint64 `json:"client_id" binding:"required"`
	CandidateID uint64 `json:"candidate_id" binding:"required"`
	Remark      string `json:"remark"`
}

type DissolveMatchValidate struct {
	ClientID uint64 `json:"client_id" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}
