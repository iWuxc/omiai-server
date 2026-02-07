package biz_omiai

// ChinaRegion 全国行政区划模型
type ChinaRegion struct {
	Code       string `json:"code" gorm:"column:code;primaryKey;size:20;comment:行政区划代码"`
	Name       string `json:"name" gorm:"column:name;size:100;not null;comment:地区名称"`
	ParentCode string `json:"parent_code" gorm:"column:parent_code;size:20;index;comment:父级行政区划代码"`
	Level      int8   `json:"level" gorm:"column:level;not null;comment:层级:1=省/直辖市,2=地级市,3=区/县"`
	Pinyin     string `json:"pinyin" gorm:"column:pinyin;size:200;index;comment:拼音（用于拼音搜索）"`
	Initial    string `json:"initial" gorm:"column:initial;size:1;index;comment:首字母（A-Z快速索引）"`
	IsHot      int8   `json:"is_hot" gorm:"column:is_hot;default:0;index;comment:是否热门城市:0=否,1=是"`
	SortOrder  int    `json:"sort_order" gorm:"column:sort_order;default:0;comment:排序权重"`
}

func (t *ChinaRegion) TableName() string {
	return "china_region"
}

// ChinaRegionInterface 地区数据层接口
type ChinaRegionInterface interface {
	// 基础查询
	GetByCode(code string) (*ChinaRegion, error)
	GetByName(name string) ([]*ChinaRegion, error)

	// 层级查询
	GetProvinces() ([]*ChinaRegion, error)                           // 获取所有省份
	GetCitiesByProvince(provinceCode string) ([]*ChinaRegion, error) // 获取某省下的城市
	GetDistrictsByCity(cityCode string) ([]*ChinaRegion, error)      // 获取某城市下的区县

	// 特殊查询
	GetHotCities() ([]*ChinaRegion, error)                  // 获取热门城市
	SearchByKeyword(keyword string) ([]*ChinaRegion, error) // 按关键词搜索
	SearchByPinyin(pinyin string) ([]*ChinaRegion, error)   // 按拼音搜索
	GetFullPath(code string) ([]*ChinaRegion, error)        // 获取完整路径（省-市-区）
}
