package biz_omiai

import (
	"context"
	"omiai-server/internal/biz"
	"time"
)

// Client 客户档案模型
type Client struct {
	ID                  uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name                string    `json:"name" gorm:"column:name;size:64;not null;comment:姓名"`
	Gender              int8      `json:"gender" gorm:"column:gender;comment:性别 1男 2女"`
	Phone               string    `json:"phone" gorm:"column:phone;size:20;index;comment:联系电话"`
	Birthday            string    `json:"birthday" gorm:"column:birthday;size:20;comment:出生年月"` // 格式 YYYY-MM
	Avatar              string    `json:"avatar" gorm:"column:avatar;size:255;comment:头像URL"`
	Age                 int       `json:"age" gorm:"column:age;comment:年龄"`
	Zodiac              string    `json:"zodiac" gorm:"column:zodiac;size:10;comment:属相"`
	Height              int       `json:"height" gorm:"column:height;comment:身高cm"`
	Weight              int       `json:"weight" gorm:"column:weight;comment:体重kg"`
	Education           int8      `json:"education" gorm:"column:education;comment:学历"` // 枚举值
	MaritalStatus       int8      `json:"marital_status" gorm:"column:marital_status;comment:婚姻状况 1未婚 2已婚 3离异 4丧偶"`
	Address             string    `json:"address" gorm:"column:address;size:255;comment:家庭住址"`
	FamilyDescription   string    `json:"family_description" gorm:"column:family_description;type:text;comment:家庭成员描述"`
	Income              int       `json:"income" gorm:"column:income;comment:月收入"`
	Profession          string    `json:"profession" gorm:"column:profession;size:128;comment:具体工作"`
	HouseStatus         int8      `json:"house_status" gorm:"column:house_status;comment:房产情况 1无房 2已购房 3贷款购房"`
	HouseAddress        string    `json:"house_address" gorm:"column:house_address;size:255;comment:买房地址"`
	CarStatus           int8      `json:"car_status" gorm:"column:car_status;comment:车辆情况 1无车 2有车"`
	Status              int8      `json:"status" gorm:"column:status;default:1;comment:状态 1单身 2匹配中 3已匹配 4停止服务"`
	ManagerID           uint64    `json:"manager_id" gorm:"column:manager_id;index;default:0;comment:归属红娘ID;-"`
	IsPublic            bool      `json:"is_public" gorm:"column:is_public;default:true;index;comment:是否公海;-"`
	Tags                string    `json:"tags" gorm:"column:tags;type:text;comment:标签列表(JSON);-"`
	PartnerRequirements string    `json:"partner_requirements" gorm:"column:partner_requirements;type:text;comment:对另一半要求(JSON)"`
	Remark              string    `json:"remark" gorm:"column:remark;type:text;comment:红娘备注"`
	Photos              string    `json:"photos" gorm:"column:photos;type:text;comment:照片URL列表(JSON)"`
	CreatedAt           time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 表名
func (t *Client) TableName() string {
	return "client"
}

// ClientInterface 定义数据层接口
type ClientInterface interface {
	Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*Client, error)
	Create(ctx context.Context, client *Client) error
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*Client, error)
	Stats(ctx context.Context) (map[string]int64, error)
}
