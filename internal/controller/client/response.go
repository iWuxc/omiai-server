package client

import "time"

type ClientResponse struct {
	ID                  uint64    `json:"id"`
	Name                string    `json:"name"`
	Gender              int8      `json:"gender"`
	Phone               string    `json:"phone"`
	Birthday            string    `json:"birthday"`
	Age                 int       `json:"age"`
	Avatar              string    `json:"avatar"`
	Zodiac              string    `json:"zodiac"`
	Height              int       `json:"height"`
	Weight              int       `json:"weight"`
	Education           int8      `json:"education"`
	MaritalStatus       int8      `json:"marital_status"`
	Address             string    `json:"address"`
	FamilyDescription   string    `json:"family_description"`
	Income              int       `json:"income"`
	Profession          string    `json:"profession"`
	WorkUnit            string    `json:"work_unit"`
	WorkCity            string    `json:"work_city"`
	WorkProvinceCode    string    `json:"work_province_code"`
	WorkCityCode        string    `json:"work_city_code"`
	WorkDistrictCode    string    `json:"work_district_code"`
	Position            string    `json:"position"`
	ParentsProfession   string    `json:"parents_profession"`
	HouseStatus         int8      `json:"house_status"`
	HouseAddress        string    `json:"house_address"`
	HouseProvinceCode   string    `json:"house_province_code"`
	HouseCityCode       string    `json:"house_city_code"`
	HouseDistrictCode   string    `json:"house_district_code"`
	CarStatus           int8      `json:"car_status"`
	Status              int8      `json:"status"`
	PartnerID           uint64    `json:"partner_id"`
	PartnerName         string    `json:"partner_name,omitempty"`
	PartnerAvatar       string    `json:"partner_avatar,omitempty"`
	PartnerRequirements string    `json:"partner_requirements"`
	Remark              string    `json:"remark"`
	Photos              string    `json:"photos"`
	ManagerID           uint64    `json:"manager_id"`
	IsPublic            bool      `json:"is_public"`
	Tags                string    `json:"tags"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func CalculateAge(birthday string) int {
	if birthday == "" {
		return 0
	}
	birthTime, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthTime.Year()
	if now.YearDay() < birthTime.YearDay() {
		age--
	}
	return age
}
