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
	HouseStatus         int8      `json:"house_status"`
	HouseAddress        string    `json:"house_address"`
	CarStatus           int8      `json:"car_status"`
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
