package omiai

import (
	"time"

	"gorm.io/gorm"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
)

type TemplateRepo struct {
	db *data.DB
}

func NewTemplateRepo(db *data.DB) biz_omiai.TemplateRepo {
	return &TemplateRepo{db: db}
}

func (r *TemplateRepo) Create(template *biz_omiai.CommunicationTemplate) error {
	return r.db.DB.Create(template).Error
}

func (r *TemplateRepo) Update(template *biz_omiai.CommunicationTemplate) error {
	return r.db.DB.Save(template).Error
}

func (r *TemplateRepo) Delete(id int64) error {
	return r.db.DB.Delete(&biz_omiai.CommunicationTemplate{}, id).Error
}

func (r *TemplateRepo) Get(id int64) (*biz_omiai.CommunicationTemplate, error) {
	var template biz_omiai.CommunicationTemplate
	if err := r.db.DB.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepo) List(category string, page, pageSize int) ([]*biz_omiai.CommunicationTemplate, int64, error) {
	var templates []*biz_omiai.CommunicationTemplate
	var total int64

	db := r.db.DB.Model(&biz_omiai.CommunicationTemplate{})
	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("usage_count desc, created_at desc").Offset(offset).Limit(pageSize).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *TemplateRepo) IncrementUsage(id int64) error {
	return r.db.DB.Model(&biz_omiai.CommunicationTemplate{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"usage_count": gorm.Expr("usage_count + ?", 1),
		"updated_at":  time.Now(),
	}).Error
}
