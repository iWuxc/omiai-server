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
	WorkUnit            string `json:"work_unit" binding:"required"`
	Position            string `json:"position" binding:"required"`
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
	WorkUnit            string `json:"work_unit"`
	Position            string `json:"position"`
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
}

type ClientDetailValidate struct {
	ID uint64 `uri:"id" binding:"required"`
}
