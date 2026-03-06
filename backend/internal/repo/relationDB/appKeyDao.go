package relationDB

import (
	"context"

	"gorm.io/gorm"
)

type AppKeyRepo interface {
	// 插入APP密钥
	Insert(ctx context.Context, db *gorm.DB, data *AppKeyModel) error
	// 根据ID删除APP密钥
	DeleteById(ctx context.Context, db *gorm.DB, id int64) error
	// 根据ID查询APP密钥
	FindById(ctx context.Context, db *gorm.DB, id int64) (*AppKeyModel, error)
	// 根据AppSecretKey查询APP密钥
	FindByAppSecretKey(ctx context.Context, db *gorm.DB, appSecretKey string) (*AppKeyModel, error)
	// 查询所有APP密钥列表
	ListAll(ctx context.Context, db *gorm.DB) ([]*AppKeyModel, error)
	// 更新APP密钥
	Update(ctx context.Context, db *gorm.DB, data *AppKeyModel) error
	// 查询APP密钥数量
	Count(ctx context.Context, db *gorm.DB) (int64, error)
}

type AppKeyMapper struct {
}

func NewAppKeyRepo() AppKeyRepo {
	return &AppKeyMapper{}
}

// Insert 插入APP密钥
func (m *AppKeyMapper) Insert(ctx context.Context, db *gorm.DB, data *AppKeyModel) error {
	return db.WithContext(ctx).Create(data).Error
}

// DeleteById 根据ID删除APP密钥
func (m *AppKeyMapper) DeleteById(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&AppKeyModel{}, id).Error
}

// FindById 根据ID查询APP密钥
func (m *AppKeyMapper) FindById(ctx context.Context, db *gorm.DB, id int64) (*AppKeyModel, error) {
	var data AppKeyModel
	err := db.WithContext(ctx).First(&data, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

// FindByAppSecretKey 根据AppSecretKey查询APP密钥
func (m *AppKeyMapper) FindByAppSecretKey(ctx context.Context, db *gorm.DB, appSecretKey string) (*AppKeyModel, error) {
	var data AppKeyModel
	err := db.WithContext(ctx).Where("app_secret_key = ?", appSecretKey).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

// ListAll 查询所有APP密钥列表
func (m *AppKeyMapper) ListAll(ctx context.Context, db *gorm.DB) ([]*AppKeyModel, error) {
	var list []*AppKeyModel
	err := db.WithContext(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Update 更新APP密钥
func (m *AppKeyMapper) Update(ctx context.Context, db *gorm.DB, data *AppKeyModel) error {
	return db.WithContext(ctx).Save(data).Error
}

// Count 查询APP密钥数量
func (m *AppKeyMapper) Count(ctx context.Context, db *gorm.DB) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&AppKeyModel{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
