package validates

type ClientCreateValidate struct {
	Name                string `json:"name" binding:"required"`
	Gender              int8   `json:"gender" binding:"required,oneof=1 2"`
	Phone               string `json:"phone" binding:"required,len=11"`
	Birthday            string `json:"birthday" binding:"required"`
	Avatar              string `json:"avatar"`
	Age                 int    `json:"age" binding:"required,min=18"`
	Zodiac              string `json:"zodiac" binding:"required"`
	Height              int    `json:"height" binding:"required,min=100"`
	Weight              int    `json:"weight" binding:"required,min=30"`
	Education           int8   `json:"education" binding:"required"`
	MaritalStatus       int8   `json:"marital_status" binding:"required"`
	Address             string `json:"address" binding:"required"`
	FamilyDescription   string `json:"family_description" binding:"required"`
	Income              int    `json:"income" binding:"required"`
	Profession          string `json:"profession" binding:"required"`
	HouseStatus         int8   `json:"house_status" binding:"required"`
	HouseAddress        string `json:"house_address"`
	CarStatus           int8   `json:"car_status" binding:"required"`
	PartnerRequirements string `json:"partner_requirements" binding:"required"`
	Remark              string `json:"remark"`
	Photos              string `json:"photos"`
}

type ClientUpdateValidate struct {
	ID                  uint64 `json:"id" binding:"required"`
	Name                string `json:"name"`
	Gender              int8   `json:"gender"`
	Phone               string `json:"phone"`
	Birthday            string `json:"birthday"`
	Avatar              string `json:"avatar"`
	Age                 int    `json:"age"`
	Zodiac              string `json:"zodiac"`
	Height              int    `json:"height"`
	Weight              int    `json:"weight"`
	Education           int8   `json:"education"`
	MaritalStatus       int8   `json:"marital_status"`
	Address             string `json:"address"`
	FamilyDescription   string `json:"family_description"`
	Income              int    `json:"income"`
	Profession          string `json:"profession"`
	HouseStatus         int8   `json:"house_status"`
	HouseAddress        string `json:"house_address"`
	CarStatus           int8   `json:"car_status"`
	PartnerRequirements string `json:"partner_requirements"`
	Remark              string `json:"remark"`
	Photos              string `json:"photos"`
}

type ClientListValidate struct {
	Paginate
	Name       string `json:"name" form:"name"`
	Phone      string `json:"phone" form:"phone"`
	Gender     int8   `json:"gender" form:"gender"`
	MinAge     int    `json:"min_age" form:"min_age"`
	MaxAge     int    `json:"max_age" form:"max_age"`
	MinHeight  int    `json:"min_height" form:"min_height"`
	MaxHeight  int    `json:"max_height" form:"max_height"`
	MinIncome  int    `json:"min_income" form:"min_income"`
	Education  int8   `json:"education" form:"education"`
	Address    string `json:"address" form:"address"`
	Profession string `json:"profession" form:"profession"`
	// Phase 1 新增字段
	Scope         string `json:"scope" form:"scope"`                   // my | public
	Status        int8   `json:"status" form:"status"`                 // 1单身 2匹配中...
	Tags          string `json:"tags" form:"tags"`                     // 标签搜索
	IsPublic      *bool  `json:"is_public" form:"is_public"`           // 用于管理员管理
	MaritalStatus int8   `json:"marital_status" form:"marital_status"` // 婚姻状况
	HouseStatus   int8   `json:"house_status" form:"house_status"`     // 房产情况
	CarStatus     int8   `json:"car_status" form:"car_status"`         // 车辆情况
	WorkCity      string `json:"work_city" form:"work_city"`           // 工作城市
}

type ClientDetailValidate struct {
	ID uint64 `uri:"id" binding:"required"`
}
