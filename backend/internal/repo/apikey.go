package repo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	domaincommon "backend/internal/domain/common"

	"gorm.io/gorm"
)

type APIKey struct {
	ID           int64  `gorm:"column:id;type:BIGINT;primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"column:name;type:VARCHAR(100);not null;uniqueIndex:idx_sys_api_key_name_owner" json:"name"`
	KeyPrefix    string `gorm:"column:key_prefix;type:VARCHAR(32);not null;default:''" json:"keyPrefix"`
	KeySuffix    string `gorm:"column:key_suffix;type:VARCHAR(16);not null;default:''" json:"keySuffix"`
	OwnerID      int64  `gorm:"column:owner_id;type:BIGINT;not null;default:0;index;uniqueIndex:idx_sys_api_key_name_owner" json:"ownerId"`
	OwnerType    string `gorm:"column:owner_type;type:VARCHAR(32);not null;default:'personal';uniqueIndex:idx_sys_api_key_name_owner" json:"ownerType"`
	UsageType    string `gorm:"column:usage_type;type:VARCHAR(32);not null;default:'external'" json:"usageType"`
	KeyType      string `gorm:"-" json:"keyType"`
	Permission   string `gorm:"column:permission;type:VARCHAR(32);not null;default:'read_only'" json:"permission"`
	ResourceKeys string `gorm:"column:resource_keys;type:TEXT;not null;default:''" json:"resourceKeys"`
	Status       int64  `gorm:"column:status;type:BIGINT;not null;default:1;index" json:"status"`
	LastUsedTime int64  `gorm:"column:last_used_time;type:BIGINT;not null;default:0" json:"lastUsedTime"`
	SoftOnlyTime
}

func (APIKey) TableName() string { return "sys_api_key" }

type APIKeyValidation struct {
	ID           int64
	OwnerID      int64
	ResourceKeys []string
	Name         string
	KeyPrefix    string
	KeyType      string
	Permission   string
}

// sysApiKey carries the key_hash column that never leaves the repo layer.
type sysApiKey struct {
	ID           int64  `gorm:"column:id;type:BIGINT;primaryKey;autoIncrement"`
	Name         string `gorm:"column:name;type:VARCHAR(100);not null;uniqueIndex:idx_sys_api_key_name_owner"`
	KeyHash      string `gorm:"column:key_hash;type:VARCHAR(128);not null;uniqueIndex:idx_sys_api_key_hash"`
	KeyPrefix    string `gorm:"column:key_prefix;type:VARCHAR(32);not null;default:''"`
	KeySuffix    string `gorm:"column:key_suffix;type:VARCHAR(16);not null;default:''"`
	OwnerID      int64  `gorm:"column:owner_id;type:BIGINT;not null;default:0;index;uniqueIndex:idx_sys_api_key_name_owner"`
	OwnerType    string `gorm:"column:owner_type;type:VARCHAR(32);not null;default:'personal';uniqueIndex:idx_sys_api_key_name_owner"`
	UsageType    string `gorm:"column:usage_type;type:VARCHAR(32);not null;default:'external'"`
	Permission   string `gorm:"column:permission;type:VARCHAR(32);not null;default:'read_only'"`
	ResourceKeys string `gorm:"column:resource_keys;type:TEXT;not null;default:''"`
	Status       int64  `gorm:"column:status;type:BIGINT;not null;default:1;index"`
	LastUsedTime int64  `gorm:"column:last_used_time;type:BIGINT;not null;default:0"`
	SoftOnlyTime
}

func (sysApiKey) TableName() string { return "sys_api_key" }

func (r *APIKeyRepo) ListAPIKeys(ctx context.Context, ownerID int64, keyType, keyword string) ([]APIKey, error) {
	var out []APIKey
	query := r.db.WithContext(ctx).
		Where("owner_id = ? AND deleted_time = 0", ownerID)
	if strings.TrimSpace(keyType) != "" {
		query = query.Where("owner_type = ?", normalizeAPIKeyOwnerType(keyType))
	}
	if strings.TrimSpace(keyword) != "" {
		query = query.Where("name ILIKE ?", "%"+strings.TrimSpace(keyword)+"%")
	}
	query = query.Where("usage_type != ? OR usage_type IS NULL", "internal")
	err := query.Order("id DESC").Find(&out).Error
	for i := range out {
		out[i].KeyType = apiKeyOwnerType(out[i].OwnerType)
	}
	return out, err
}

func (r *APIKeyRepo) CreateAPIKey(ctx context.Context, name string, ownerID, workspaceID int64, permission, usageType, keyType string, resourceKeys []string) (APIKey, string, error) {
	if usageType == "" {
		usageType = "external"
	}
	if permission == "" {
		permission = "read_only"
	}
	if keyType == "" {
		keyType = "personal"
	}
	ownerType := normalizeAPIKeyOwnerType(keyType)
	raw, err := GenerateWorkspaceAPIKey(ownerType, workspaceID)
	if err != nil {
		return APIKey{}, "", err
	}
	prefix, suffix := secretEdges(raw)
	row := sysApiKey{
		Name:         name,
		KeyHash:      hashSecret(raw),
		KeyPrefix:    prefix,
		KeySuffix:    suffix,
		OwnerID:      ownerID,
		OwnerType:    ownerType,
		UsageType:    usageType,
		Permission:   permission,
		ResourceKeys: strings.Join(resourceKeys, ","),
		Status:       1,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return APIKey{}, "", normalizeDBError(err)
	}
	item := APIKey{
		ID:           row.ID,
		Name:         row.Name,
		KeyPrefix:    row.KeyPrefix,
		KeySuffix:    row.KeySuffix,
		OwnerID:      row.OwnerID,
		OwnerType:    row.OwnerType,
		UsageType:    row.UsageType,
		KeyType:      apiKeyOwnerType(row.OwnerType),
		Permission:   row.Permission,
		ResourceKeys: row.ResourceKeys,
		Status:       row.Status,
		LastUsedTime: row.LastUsedTime,
	}
	item.CreatedTime = row.CreatedTime
	item.UpdatedTime = row.UpdatedTime
	item.DeletedTime = row.DeletedTime
	return item, raw, nil
}

// GenerateWorkspaceAPIKey mirrors Cloud API key encoding so the shared SDK can
// derive the MQTT username and client ID without edition-specific branches.
func GenerateWorkspaceAPIKey(keyType string, workspaceID int64) (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	prefix := "sk-per-"
	switch normalizeAPIKeyOwnerType(keyType) {
	case "service":
		prefix = "sk-svc-"
	case "agent":
		prefix = "sk-agent-"
	}
	raw := domaincommon.FormatWorkspaceAPIKey(prefix, workspaceID, base64.RawURLEncoding.EncodeToString(token))
	if raw == "" {
		return "", ErrInvalidArgument
	}
	return raw, nil
}

// RegisterAPIKey 用外部已生成的 raw key 注册到 sys_api_key（供 OpenAPI 鉴权）。
// 与 CreateAPIKey 的区别：raw 由调用方提供，且不返回明文（调用方已持有）。
func (r *APIKeyRepo) RegisterAPIKey(ctx context.Context, raw, name string, ownerID int64, keyType, permission, usageType string, resourceKeys []string) error {
	raw = strings.TrimSpace(raw)
	name = strings.TrimSpace(name)
	if raw == "" || name == "" || ownerID <= 0 {
		return ErrInvalidArgument
	}
	if usageType == "" {
		usageType = "external"
	}
	permission = NormalizeAPIKeyPermission(permission)
	ownerType := normalizeAPIKeyOwnerType(keyType)
	keyHash := hashSecret(raw)
	prefix, suffix := secretEdges(raw)

	updates := map[string]any{
		"name":          name,
		"key_hash":      keyHash,
		"key_prefix":    prefix,
		"key_suffix":    suffix,
		"owner_id":      ownerID,
		"owner_type":    ownerType,
		"usage_type":    usageType,
		"permission":    permission,
		"resource_keys": strings.Join(resourceKeys, ","),
		"status":        1,
		"deleted_time":  0,
	}

	var existing sysApiKey
	err := r.db.WithContext(ctx).Unscoped().Where("key_hash = ?", keyHash).Take(&existing).Error
	if err == nil {
		if existing.DeletedTime == 0 {
			return nil
		}
		return normalizeDBError(r.db.WithContext(ctx).Unscoped().Model(&sysApiKey{}).
			Where("id = ?", existing.ID).
			Updates(touchValues(updates, 0)).Error)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return normalizeDBError(err)
	}

	err = r.db.WithContext(ctx).
		Where("name = ? AND owner_id = ? AND owner_type = ? AND deleted_time = 0", name, ownerID, ownerType).
		Take(&existing).Error
	if err == nil {
		return normalizeDBError(r.db.WithContext(ctx).Model(&sysApiKey{}).
			Where("id = ?", existing.ID).
			Updates(touchValues(updates, 0)).Error)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return normalizeDBError(err)
	}

	row := sysApiKey{
		Name:         name,
		KeyHash:      keyHash,
		KeyPrefix:    prefix,
		KeySuffix:    suffix,
		OwnerID:      ownerID,
		OwnerType:    ownerType,
		UsageType:    usageType,
		Permission:   permission,
		ResourceKeys: strings.Join(resourceKeys, ","),
		Status:       1,
	}
	return normalizeDBError(r.db.WithContext(ctx).Create(&row).Error)
}

func (r *APIKeyRepo) UpdateAPIKey(ctx context.Context, id, ownerID int64, name, permission string, resourceKeys []string) (APIKey, error) {
	updates := touchValues(map[string]any{}, 0)
	if strings.TrimSpace(name) != "" {
		updates["name"] = strings.TrimSpace(name)
	}
	if strings.TrimSpace(permission) != "" {
		updates["permission"] = strings.TrimSpace(permission)
		updates["resource_keys"] = strings.Join(resourceKeys, ",")
	}
	res := r.db.WithContext(ctx).Model(&APIKey{}).
		Where("id = ? AND owner_id = ? AND deleted_time = 0", id, ownerID).
		Updates(updates)
	if res.Error != nil {
		return APIKey{}, normalizeDBError(res.Error)
	}
	if res.RowsAffected == 0 {
		return APIKey{}, ErrNotFound
	}
	var item APIKey
	if err := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ? AND deleted_time = 0", id, ownerID).
		Take(&item).Error; err != nil {
		return APIKey{}, err
	}
	item.KeyType = apiKeyOwnerType(item.OwnerType)
	return item, nil
}

func (r *APIKeyRepo) DeleteAPIKey(ctx context.Context, id, ownerID int64) error {
	res := r.db.WithContext(ctx).Model(&APIKey{}).
		Where("id = ? AND owner_id = ? AND deleted_time = 0", id, ownerID).
		Updates(softDeleteNoDelByValues(0, 0))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *APIKeyRepo) ValidateAPIKey(ctx context.Context, raw string) (int64, []string, error) {
	detail, err := r.ValidateAPIKeyDetail(ctx, raw)
	if err != nil {
		return 0, nil, err
	}
	return detail.OwnerID, detail.ResourceKeys, nil
}

func (r *APIKeyRepo) ValidateAPIKeyDetail(ctx context.Context, raw string) (APIKeyValidation, error) {
	return r.validateAPIKeyDetail(ctx, "key_hash = ?", hashSecret(raw))
}

func (r *APIKeyRepo) ValidateAPIKeyByIDDetail(ctx context.Context, id int64) (APIKeyValidation, error) {
	if id <= 0 {
		return APIKeyValidation{}, ErrNotFound
	}
	return r.validateAPIKeyDetail(ctx, "id = ?", id)
}

func (r *APIKeyRepo) validateAPIKeyDetail(ctx context.Context, condition string, value any) (APIKeyValidation, error) {
	var row struct {
		ID           int64  `gorm:"column:id"`
		OwnerID      int64  `gorm:"column:owner_id"`
		ResourceKeys string `gorm:"column:resource_keys"`
		Name         string `gorm:"column:name"`
		KeyPrefix    string `gorm:"column:key_prefix"`
		OwnerType    string `gorm:"column:owner_type"`
		Permission   string `gorm:"column:permission"`
	}
	if err := r.db.WithContext(ctx).Table("sys_api_key").
		Select("id, owner_id, resource_keys, name, key_prefix, owner_type, permission").
		Where(condition+" AND status = 1 AND deleted_time = 0", value).
		Take(&row).Error; err != nil {
		return APIKeyValidation{}, err
	}
	now := time.Now().UTC().UnixMilli()
	_ = r.db.WithContext(ctx).Model(&APIKey{}).
		Where("id = ? AND deleted_time = 0", row.ID).
		Update("last_used_time", now).Error
	return APIKeyValidation{
		ID:           row.ID,
		OwnerID:      row.OwnerID,
		ResourceKeys: splitCSV(row.ResourceKeys),
		Name:         row.Name,
		KeyPrefix:    row.KeyPrefix,
		KeyType:      apiKeyOwnerType(row.OwnerType),
		Permission:   row.Permission,
	}, nil
}

func normalizeAPIKeyOwnerType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "service", "workspace", "service_key":
		return "service"
	case "agent":
		return "agent"
	default:
		return "personal"
	}
}

func apiKeyOwnerType(value string) string {
	normalized := normalizeAPIKeyOwnerType(value)
	if normalized == "" {
		return "personal"
	}
	return normalized
}

type APIKeyRepo struct{ db *gorm.DB }

func NewAPIKeyRepo(in any) *APIKeyRepo { return &APIKeyRepo{db: GetCommonConn(in)} }
