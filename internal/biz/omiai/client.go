package biz_omiai

import (
	"context"
	"omiai-server/internal/biz"
	"time"
)

// Client 客户档案模型
type Client struct {
	ID                uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID          uint64 `json:"tenant_id" gorm:"column:tenant_id;index;default:1;comment:所属租户(机构)ID"`
	Name              string `json:"name" gorm:"column:name;size:64;not null;comment:姓名"`
	Gender            int8   `json:"gender" gorm:"column:gender;comment:性别 1男 2女"`
	Phone             string `json:"phone" gorm:"column:phone;size:20;index;comment:联系电话"`
	Birthday          string `json:"birthday" gorm:"column:birthday;size:20;comment:出生年月"` // 格式 YYYY-MM
	Avatar            string `json:"avatar" gorm:"column:avatar;size:255;comment:头像URL"`
	Age               int    `json:"age" gorm:"column:age;comment:年龄"`
	Zodiac            string `json:"zodiac" gorm:"column:zodiac;size:10;comment:属相"`
	Height            int    `json:"height" gorm:"column:height;comment:身高cm"`
	Weight            int    `json:"weight" gorm:"column:weight;comment:体重kg"`
	Education         int8   `json:"education" gorm:"column:education;comment:学历"` // 枚举值
	MaritalStatus     int8   `json:"marital_status" gorm:"column:marital_status;comment:婚姻状况 1未婚 2已婚 3离异 4丧偶"`
	Address           string `json:"address" gorm:"column:address;size:255;comment:家庭住址"`
	FamilyDescription string `json:"family_description" gorm:"column:family_description;type:text;comment:家庭成员描述"`
	Income            int    `json:"income" gorm:"column:income;comment:月收入"`
	Profession        string `json:"profession" gorm:"column:profession;size:128;comment:具体工作"`
	WorkUnit          string `json:"work_unit" gorm:"column:work_unit;size:128;comment:工作单位"`
	WorkCity          string `json:"work_city" gorm:"column:work_city;size:128;comment:工作城市"`
	WorkProvinceCode  string `json:"work_province_code" gorm:"column:work_province_code;size:20;comment:工作省份代码"`
	WorkCityCode      string `json:"work_city_code" gorm:"column:work_city_code;size:20;comment:工作城市代码"`
	WorkDistrictCode  string `json:"work_district_code" gorm:"column:work_district_code;size:20;comment:工作区县代码"`

	Position            string    `json:"position" gorm:"column:position;size:128;comment:职位"`
	HouseStatus         int8      `json:"house_status" gorm:"column:house_status;comment:房产情况 1无房 2已购房 3贷款购房"`
	HouseAddress        string    `json:"house_address" gorm:"column:house_address;size:255;comment:买房地址"`
	HouseProvinceCode   string    `json:"house_province_code" gorm:"column:house_province_code;size:20;comment:房产省份代码"`
	HouseCityCode       string    `json:"house_city_code" gorm:"column:house_city_code;size:20;comment:房产城市代码"`
	HouseDistrictCode   string    `json:"house_district_code" gorm:"column:house_district_code;size:20;comment:房产区县代码"`
	CarStatus           int8      `json:"car_status" gorm:"column:car_status;comment:车辆情况 1无车 2有车"`
	Status              int8      `json:"status" gorm:"column:status;default:1;comment:状态 1单身 2匹配中 3已匹配 4停止服务"`
	PartnerID           *uint64   `json:"partner_id" gorm:"column:partner_id;uniqueIndex;default:null;comment:当前匹配对象ID"`
	Partner             *Client   `json:"partner" gorm:"foreignKey:PartnerID"`
	ManagerID           uint64    `json:"manager_id" gorm:"column:manager_id;index;default:0;comment:归属红娘ID;-"`
	IsPublic            bool      `json:"is_public" gorm:"column:is_public;default:true;index;comment:是否公海;-"`
	Tags                string    `json:"tags" gorm:"column:tags;type:text;comment:系统兴趣标签列表(JSON);-"`
	InterestTags        string    `json:"interest_tags" gorm:"column:interest_tags;type:text;comment:用户自定义兴趣标签(JSON)"`
	PartnerRequirements string    `json:"partner_requirements" gorm:"column:partner_requirements;type:text;comment:对另一半要求(JSON)"`
	ParentsProfession   string    `json:"parents_profession" gorm:"column:parents_profession;size:255;comment:父母工作"`
	Remark              string    `json:"remark" gorm:"column:remark;type:text;comment:红娘备注"`
	Photos              string    `json:"photos" gorm:"column:photos;type:text;comment:照片URL列表(JSON)`
	CandidateCacheJSON  string    `json:"candidate_cache_json" gorm:"column:candidate_cache_json;type:text;comment:算法初筛结果缓存"`
	WxOpenid            string    `json:"wx_openid" gorm:"column:wx_openid;size:128;uniqueIndex;comment:微信OpenID"`
	WxUnionid           string    `json:"wx_unionid" gorm:"column:wx_unionid;size:128;comment:微信UnionID"`
	IsVerified          bool      `json:"is_verified" gorm:"column:is_verified;default:false;comment:是否已红娘审核"`
	IdCardNo            string    `json:"id_card_no" gorm:"column:id_card_no;size:18;comment:身份证号(加密)"`
	RealName            string    `json:"real_name" gorm:"column:real_name;size:64;comment:真实姓名"`
	IsRealNameVerified  bool      `json:"is_real_name_verified" gorm:"column:is_real_name_verified;default:false;comment:是否已实名认证"`
	Coins               int       `json:"coins" gorm:"column:coins;default:0;comment:虚拟币余额(红豆)"`
	VipExpireAt         time.Time `json:"vip_expire_at" gorm:"column:vip_expire_at;comment:VIP到期时间"`
	CreatedAt           time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// ClientCoinRecord 虚拟币流水表
type ClientCoinRecord struct {
	ID        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientID  uint64    `json:"client_id" gorm:"column:client_id;index;comment:客户ID"`
	Amount    int       `json:"amount" gorm:"column:amount;comment:变动金额(正负)"`
	Type      int8      `json:"type" gorm:"column:type;comment:变动类型 1充值 2解锁查看 3开通VIP"`
	Remark    string    `json:"remark" gorm:"column:remark;size:255;comment:备注说明"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (t *ClientCoinRecord) TableName() string {
	return "client_coin_record"
}

// ClientInteraction 互动记录表
type ClientInteraction struct {
	ID           uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID     uint64    `json:"tenant_id" gorm:"column:tenant_id;index;default:1;comment:所属租户(机构)ID"`
	FromClientID uint64    `json:"from_client_id" gorm:"column:from_client_id;index;comment:发起方ID"`
	ToClientID   uint64    `json:"to_client_id" gorm:"column:to_client_id;index;comment:接收方ID"`
	ActionType   int8      `json:"action_type" gorm:"column:action_type;comment:行为类型 1查看 2单向心动 3互相心动"`
	Status       int8      `json:"status" gorm:"column:status;default:0;comment:状态 0未处理 1已跟进 2已忽略"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (t *ClientInteraction) TableName() string {
	return "client_interaction"
}

// TableName 表名
func (t *Client) TableName() string {
	return "client"
}

func (c *Client) RealAge() int {
	if c.Age > 0 {
		return c.Age
	}
	if len(c.Birthday) < 7 {
		return 0
	}
	t, err := time.Parse("2006-01", c.Birthday[:7])
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - t.Year()
	if now.Month() < t.Month() {
		age--
	}
	return age
}

// ClientInterface 定义数据层接口
type ClientInterface interface {
	Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*Client, error)
	Create(ctx context.Context, client *Client) error
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*Client, error)
	Stats(ctx context.Context) (map[string]int64, error)

	// C端相关
	GetByWxOpenID(ctx context.Context, openID string) (*Client, error)
	SaveInteraction(ctx context.Context, interaction *ClientInteraction) error
	GetInteraction(ctx context.Context, fromID, toID uint64) (*ClientInteraction, error)
	GetInteractionLeads(ctx context.Context, managerID uint64, offset, limit int) ([]*ClientInteraction, error)
	GetClientInteractions(ctx context.Context, clientID uint64, actionType int8, offset, limit int) ([]*ClientInteraction, error)

	// 商业化相关
	AddCoins(ctx context.Context, clientID uint64, amount int, recordType int8, remark string) error
	IsVip(ctx context.Context, clientID uint64) bool

	// Dashboard 相关
	GetDashboardStats(ctx context.Context) (map[string]int64, error)
}
