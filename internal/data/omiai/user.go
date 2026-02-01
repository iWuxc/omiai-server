package omiai

import (
	"context"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"gorm.io/gorm"
)

var _ biz_omiai.UserInterface = (*UserRepo)(nil)

type UserRepo struct {
	db *data.DB
}

func NewUserRepo(db *data.DB) biz_omiai.UserInterface {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*biz_omiai.User, error) {
	var user biz_omiai.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByWxOpenID(ctx context.Context, openID string) (*biz_omiai.User, error) {
	var user biz_omiai.User
	err := r.db.WithContext(ctx).Where("wx_openid = ?", openID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *biz_omiai.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) Update(ctx context.Context, user *biz_omiai.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*biz_omiai.User, error) {
	var user biz_omiai.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
