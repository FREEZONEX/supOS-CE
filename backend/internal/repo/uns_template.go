package repo

import (
	"context"
	"encoding/json"
)

type UnsTemplate struct {
	ID               int64           `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	Name             string          `gorm:"column:name" json:"name"`
	TopicType        int16           `gorm:"column:topic_type" json:"topicType"`
	Schema           json.RawMessage `gorm:"column:schema;type:jsonb" json:"schema"`
	ExtendProperties json.RawMessage `gorm:"column:extend_properties;type:jsonb" json:"extendProperties"`
	SoftTime
}

func (UnsTemplate) TableName() string { return "uns_namespace_template_info" }

func (r *UnsRepo) ListUnsTemplates(ctx context.Context) ([]UnsTemplate, error) {
	var out []UnsTemplate
	err := r.db.WithContext(ctx).
		Select("id, name, topic_type, schema, extend_properties").
		Where("deleted_time = 0").Order("id").Find(&out).Error
	return out, err
}

func (r *UnsRepo) GetUnsTemplate(ctx context.Context, id int64) (UnsTemplate, error) {
	var item UnsTemplate
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_time = 0", id).Take(&item).Error
	return item, err
}

func (r *UnsRepo) CreateUnsTemplate(ctx context.Context, item UnsTemplate) (UnsTemplate, error) {
	ensureID(&item.ID)
	if len(item.Schema) == 0 {
		item.Schema = json.RawMessage(`{}`)
	}
	if len(item.ExtendProperties) == 0 {
		item.ExtendProperties = json.RawMessage(`{}`)
	}
	if item.TopicType == 0 {
		item.TopicType = 1
	}
	err := r.db.WithContext(ctx).Create(&item).Error
	return item, normalizeDBError(err)
}

func (r *UnsRepo) UpdateUnsTemplate(ctx context.Context, id int64, item UnsTemplate) (UnsTemplate, error) {
	now := repoNowMilli()
	if len(item.Schema) == 0 {
		item.Schema = json.RawMessage(`{}`)
	}
	if len(item.ExtendProperties) == 0 {
		item.ExtendProperties = json.RawMessage(`{}`)
	}
	if item.TopicType == 0 {
		item.TopicType = 1
	}
	var out UnsTemplate
	res := r.db.WithContext(ctx).Model(&out).Clauses(returningAll()).
		Where("id = ? AND deleted_time = 0", id).
		Updates(touchByValues(map[string]any{
			"name":              item.Name,
			"topic_type":        item.TopicType,
			"schema":            item.Schema,
			"extend_properties": item.ExtendProperties,
		}, item.UpdatedBy, now))
	if res.Error != nil {
		return UnsTemplate{}, normalizeDBError(res.Error)
	}
	if res.RowsAffected == 0 {
		return UnsTemplate{}, ErrNotFound
	}
	return out, nil
}

func (r *UnsRepo) DeleteUnsTemplate(ctx context.Context, id, userID int64) error {
	res := r.db.WithContext(ctx).Model(&UnsTemplate{}).
		Where("id = ? AND deleted_time = 0", id).
		Updates(softDeleteValues(userID, 0))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
