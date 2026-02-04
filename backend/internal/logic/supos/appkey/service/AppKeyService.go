package service

import (
	"backend/internal/common"
	"backend/internal/config"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
}

const (
	StatusEnabled  int32 = 1 // 启用
	StatusDisabled int32 = 0 // 禁用
)

type KongaSecretKeyVo struct {
	ID string `json:"id"`
}

type AppKeyService struct {
	appKeyRepo dao.AppKeyRepo
	host       string
	port       string
	path       string
}

func init() {
	spring.RegisterBean[*AppKeyService](&AppKeyService{
		appKeyRepo: dao.NewAppKeyRepo(), // 初始化 appKeyRepo
	})
}

func (s *AppKeyService) Init(c config.Config) error {
	// 初始化Konga配置，这里可以从配置文件读取
	s.host = c.Kong.Host
	s.port = strconv.Itoa(c.Kong.Port)
	s.path = "kong/home/consumers/59d1ef15-24a5-4373-b957-e8192c15ff6e/key-auth"
	return nil
}

// CreateSecretKey 创建密钥
func (s *AppKeyService) CreateSecretKey(ctx context.Context) (bool, error) {
	db := dao.GetDb(ctx)

	count, err := s.appKeyRepo.Count(ctx, db)
	if err != nil {
		return false, err
	}

	if count >= 3 {
		return false, errors.NewCodeError(400, "最多只能新增3个密钥")
	}

	key := uuid.NewString()

	// 调用Konga API创建密钥
	kongaKey, err := s.createKongaKey(key)
	if err != nil {
		return false, err
	}

	// 保存到数据库
	appKey := &dao.AppKeyModel{
		ID:             common.NextId(),
		AppSecretKey:   key,
		AppSecretValue: kongaKey.ID,
		Status:         StatusEnabled,
		CreateTime:     time.Now(),
	}

	err = s.appKeyRepo.Insert(ctx, db, appKey)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateSecretKey 更新密钥状态
func (s *AppKeyService) UpdateSecretKey(ctx context.Context, req *types.UpdateAppKeyReq) error {
	db := dao.GetDb(ctx)

	appKey, err := s.appKeyRepo.FindByAppSecretKey(ctx, db, req.AppSecretKey)
	if err != nil {
		return err
	}

	if appKey == nil {
		return errors.NewCodeError(400, "密钥不存在")
	}

	// 开始事务
	err = db.Transaction(func(tx *gorm.DB) error {
		if req.Status == StatusEnabled {
			// 启用密钥，调用Konga API生成新的secretValue
			kongaKey, err := s.createKongaKey(appKey.AppSecretKey)
			if err != nil {
				return err
			}
			appKey.AppSecretValue = kongaKey.ID
		} else {
			// 禁用密钥，删除Konga中的记录
			err := s.deleteKongaKey(appKey.AppSecretValue)
			if err != nil {
				return err
			}
		}

		appKey.Status = req.Status
		return s.appKeyRepo.Update(ctx, tx, appKey)
	})

	return err
}

// GetSecretKeyList 获取密钥列表
func (s *AppKeyService) GetSecretKeyList(ctx context.Context) ([]*types.AppKeyInfo, error) {
	db := dao.GetDb(ctx)

	list, err := s.appKeyRepo.ListAll(ctx, db)

	if err != nil {
		return nil, err
	}

	resp := make([]*types.AppKeyInfo, 0, len(list))
	for _, item := range list {
		resp = append(resp, &types.AppKeyInfo{
			ID:             strconv.FormatInt(item.ID, 10),
			AppSecretKey:   item.AppSecretKey,
			AppSecretValue: item.AppSecretValue,
			Status:         item.Status,
			CreateTime:     item.CreateTime.UnixMilli(),
		})
	}

	return resp, nil
}

// DeleteSecretKey 删除密钥
func (s *AppKeyService) DeleteSecretKey(ctx context.Context, id int64) error {
	db := dao.GetDb(ctx)

	appKey, err := s.appKeyRepo.FindById(ctx, db, id)
	if err != nil {
		return err
	}

	if appKey == nil {
		return errors.NewCodeError(400, "密钥不存在")
	}

	if appKey.Status == StatusEnabled {
		return errors.NewCodeError(400, "只能删除禁用状态的密钥")
	}

	err = s.appKeyRepo.DeleteById(ctx, db, id)
	if err != nil {
		return err
	}

	return nil
}

// 构建Konga API URL
func (s *AppKeyService) buildKongaUrl(suffix string) string {
	return fmt.Sprintf("http://%s:%s/%s%s", s.host, s.port, s.path, suffix)
}

// 创建Konga密钥
func (s *AppKeyService) createKongaKey(key string) (*KongaSecretKeyVo, error) {
	url := s.buildKongaUrl("")
	body := map[string]string{
		"key": key,
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		logx.Errorf("Failed to create key in Konga: %v", err)
		return nil, errors.NewCodeError(500, "创建密钥失败")
	}

	if resp.StatusCode()/100 != 2 {
		logx.Errorf("Failed to create key in Konga: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, errors.NewCodeError(500, "创建密钥失败")
	}

	var kongaKey KongaSecretKeyVo
	err = json.Unmarshal(resp.Body(), &kongaKey)
	if err != nil {
		logx.Errorf("Failed to parse Konga response: %v, body: %s", err, resp.Body())
		return nil, errors.NewCodeError(500, "解析响应失败")
	}

	return &kongaKey, nil
}

// 删除Konga密钥
func (s *AppKeyService) deleteKongaKey(keyId string) error {
	url := s.buildKongaUrl(fmt.Sprintf("/%s", keyId))

	resp, err := resty.New().R().
		Delete(url)

	if err != nil {
		logx.Errorf("Failed to delete key from Konga: %v", err)
		return errors.NewCodeError(500, "删除密钥失败")
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 404 {
		logx.Errorf("Failed to delete key from Konga: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return errors.NewCodeError(500, "删除密钥失败")
	}

	return nil
}

// 获取Konga密钥
func (s *AppKeyService) getKongaKey(keyId string) (*KongaSecretKeyVo, error) {
	url := s.buildKongaUrl(fmt.Sprintf("/%s", keyId))

	resp, err := resty.New().R().
		Get(url)

	if err != nil {
		logx.Errorf("Failed to get key from Konga: %v", err)
		return nil, errors.NewCodeError(500, "获取密钥失败")
	}

	if resp.StatusCode() == 404 {
		return nil, nil
	}

	if resp.StatusCode()/100 != 2 {
		logx.Errorf("Failed to get key from Konga: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, errors.NewCodeError(500, "获取密钥失败")
	}

	var kongaKey KongaSecretKeyVo
	err = json.Unmarshal(resp.Body(), &kongaKey)
	if err != nil {
		logx.Errorf("Failed to parse Konga response: %v, body: %s", err, resp.Body())
		return nil, errors.NewCodeError(500, "解析响应失败")
	}

	return &kongaKey, nil
}
