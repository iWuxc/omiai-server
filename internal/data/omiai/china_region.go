package omiai

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"gorm.io/gorm"
)

type ChinaRegionRepo struct {
	db *data.DB
}

func NewChinaRegionRepo(db *data.DB) biz_omiai.ChinaRegionInterface {
	return &ChinaRegionRepo{db: db}
}

// GetByCode 根据行政区划代码获取详情
func (r *ChinaRegionRepo) GetByCode(code string) (*biz_omiai.ChinaRegion, error) {
	var region biz_omiai.ChinaRegion
	if err := r.db.DB.Where("code = ?", code).First(&region).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

// GetByName 根据名称精确查询
func (r *ChinaRegionRepo) GetByName(name string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	if err := r.db.DB.Where("name = ?", name).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetProvinces 获取所有省份 (Level=1)
func (r *ChinaRegionRepo) GetProvinces() ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	if err := r.db.DB.Where("level = ?", 1).Order("code asc").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetCitiesByProvince 获取某省下的城市 (Level=2, ParentCode=provinceCode)
func (r *ChinaRegionRepo) GetCitiesByProvince(provinceCode string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	if err := r.db.DB.Where("parent_code = ? AND level = ?", provinceCode, 2).Order("code asc").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetDistrictsByCity 获取某城市下的区县 (Level=3, ParentCode=cityCode)
func (r *ChinaRegionRepo) GetDistrictsByCity(cityCode string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	if err := r.db.DB.Where("parent_code = ? AND level = ?", cityCode, 3).Order("code asc").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetHotCities 获取热门城市
func (r *ChinaRegionRepo) GetHotCities() ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	// 通常热门城市是Level=2 (市级)，但也可能是Level=1 (直辖市)
	if err := r.db.DB.Where("is_hot = ?", 1).Order("sort_order asc, code asc").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// SearchByKeyword 按关键词模糊搜索 (Name)
func (r *ChinaRegionRepo) SearchByKeyword(keyword string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	likeStr := "%" + keyword + "%"
	if err := r.db.DB.Where("name LIKE ?", likeStr).Order("level asc, code asc").Limit(20).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// SearchByPinyin 按拼音模糊搜索 (Pinyin or Initial)
func (r *ChinaRegionRepo) SearchByPinyin(pinyin string) ([]*biz_omiai.ChinaRegion, error) {
	var regions []*biz_omiai.ChinaRegion
	likeStr := "%" + pinyin + "%"
	if err := r.db.DB.Where("pinyin LIKE ? OR initial = ?", likeStr, pinyin).Order("level asc, code asc").Limit(20).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetFullPath 获取完整路径（递归查找父级）
// 注意：MySQL 8.0+ 支持 CTE 递归查询，这里为了兼容性使用代码递归或多次查询
// 由于层级固定为3级，直接多次查询更简单
func (r *ChinaRegionRepo) GetFullPath(code string) ([]*biz_omiai.ChinaRegion, error) {
	var current biz_omiai.ChinaRegion
	if err := r.db.DB.Where("code = ?", code).First(&current).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	path := []*biz_omiai.ChinaRegion{&current}

	// 向上查找父级
	parentCode := current.ParentCode
	for parentCode != "" {
		var parent biz_omiai.ChinaRegion
		if err := r.db.DB.Where("code = ?", parentCode).First(&parent).Error; err != nil {
			break
		}
		// 插入到切片头部
		path = append([]*biz_omiai.ChinaRegion{&parent}, path...)
		parentCode = parent.ParentCode
	}

	return path, nil
}
