package client

import "time"

type ClientResponse struct {
	ID                  uint64    `json:"id"`
	Name                string    `json:"name"`
	Gender              int8      `json:"gender"`
	Phone               string    `json:"phone"`
	Birthday            string    `json:"birthday"`
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
	CarStatus           int8      `json:"car_status"`
	PartnerRequirements string    `json:"partner_requirements"`
	Remark              string    `json:"remark"`
	Photos              string    `json:"photos"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
