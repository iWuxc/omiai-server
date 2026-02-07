package omiai

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"gorm.io/gorm"
)

var _ biz_omiai.ChinaRegionInterface = (*ChinaRegionRepo)(nil)

type ChinaRegionRepo struct {
	db *data.DB
}

func NewChinaRegionRepo(db *data.DB) biz_omiai.ChinaRegionInterface {
	return &ChinaRegionRepo{db: db}
}

// GetByCode 根据编码查询
func (r *ChinaRegionRepo) GetByCode(code string) (*biz_omiai.ChinaRegion, error) {
	var region biz_omiai.ChinaRegion
	err := r.db.Where("code = ?", code).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetByName 根据名称模糊查询
func (r *ChinaRegionRepo) GetByName(name string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("name LIKE ?", "%"+name+"%").
		Order("level ASC, code ASC").
		Find(&regions).Error
	return regions, err
}

// GetProvinces 获取所有省份（level=1）
func (r *ChinaRegionRepo) GetProvinces() ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("level = ?", 1).
		Order("code ASC").
		Find(&regions).Error
	return regions, err
}

// GetCitiesByProvince 获取某省下的所有城市
func (r *ChinaRegionRepo) GetCitiesByProvince(provinceCode string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("parent_code = ? AND level = ?", provinceCode, 2).
		Order("code ASC").
		Find(&regions).Error
	return regions, err
}

// GetDistrictsByCity 获取某城市下的所有区县
func (r *ChinaRegionRepo) GetDistrictsByCity(cityCode string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("parent_code = ? AND level = ?", cityCode, 3).
		Order("code ASC").
		Find(&regions).Error
	return regions, err
}

// GetHotCities 获取热门城市（用于快捷选择）
func (r *ChinaRegionRepo) GetHotCities() ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("is_hot = ? AND level = ?", 1, 2).
		Order("sort_order ASC, code ASC").
		Find(&regions).Error
	return regions, err
}

// SearchByKeyword 关键词搜索（支持名称和拼音）
func (r *ChinaRegionRepo) SearchByKeyword(keyword string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("name LIKE ? OR pinyin LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("level ASC, code ASC").
		Limit(20).
		Find(&regions).Error
	return regions, err
}

// SearchByPinyin 拼音首字母搜索
func (r *ChinaRegionRepo) SearchByPinyin(pinyin string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	// 支持首字母和完整拼音搜索
	err := r.db.Where("initial = ? OR pinyin LIKE ?",
		pinyin, "%"+pinyin+"%").
		Order("level ASC, code ASC").
		Limit(20).
		Find(&regions).Error
	return regions, err
}

// GetFullPath 获取完整地区路径（省-市-区）
func (r *ChinaRegionRepo) GetFullPath(code string) ([]*biz_omiai.ChinaRegion, error) {
	var path []*biz_omiai.ChinaRegion

	currentCode := code
	for currentCode != "" {
		var region biz_omiai.ChinaRegion
		err := r.db.Where("code = ?", currentCode).First(&region).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return nil, err
		}

		// 插入到头部，保持省-市-区的顺序
		path = append([]*biz_omiai.ChinaRegion{&region}, path...)

		// 继续查询父级
		currentCode = region.ParentCode
	}

	return path, nil
}

// GetCitiesByInitial 按首字母获取城市（用于A-Z索引）
func (r *ChinaRegionRepo) GetCitiesByInitial(initial string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	err := r.db.Where("initial = ? AND level = ?", initial, 2).
		Order("code ASC").
		Find(&regions).Error
	return regions, err
}

// GetAllCities 获取所有城市（排除省直辖县）
func (r *ChinaRegionRepo) GetAllCities() ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	// level=2 且不是省直辖县（parent_code不以00结尾的）
	err := r.db.Where("level = ?", 2).
		Order("parent_code ASC, code ASC").
		Find(&regions).Error
	return regions, err
}
