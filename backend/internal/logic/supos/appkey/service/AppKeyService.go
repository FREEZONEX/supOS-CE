package service

import (
	"backend/internal/common"
	"backend/internal/common/event"
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

type KongSecretKeyVo struct {
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
	// 初始化 Kong Admin 配置，这里可以从配置文件读取
	s.host = c.Kong.Host
	s.port = strconv.Itoa(c.Kong.Port)
	s.path = "consumers/59d1ef15-24a5-4373-b957-e8192c15ff6e/key-auth"
	return nil
}

// OnEventContextRefreshed 监听应用上下文刷新事件，自动创建appKey
func (s *AppKeyService) OnEventContextRefreshed(event *event.ContextRefreshedEvent) error {
	logx.Infof("AppKeyService: Checking for existing secret keys...")

	// 获取密钥列表
	ctx := context.Background()
	keyList, err := s.GetSecretKeyList(ctx)
	if err != nil {
		logx.Errorf("AppKeyService: Failed to get secret key list: %v", err)
		return err
	}

	// 判断列表是否为空
	if len(keyList) == 0 {
		logx.Infof("AppKeyService: No secret keys found, automatically creating one...")
		success, err := s.CreateSecretKey(ctx)
		if err != nil {
			logx.Errorf("AppKeyService: Failed to create secret key: %v", err)
			return err
		}
		if success {
			logx.Infof("AppKeyService: Secret key created successfully")
		}
	} else {
		logx.Infof("AppKeyService: Found %d existing secret keys, no action needed", len(keyList))
	}

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

	// 调用 Kong Admin API 创建密钥
	kongKey, err := s.createKongKey(key)
	if err != nil {
		return false, err
	}

	// 保存到数据库
	appKey := &dao.AppKeyModel{
		ID:             common.NextId(),
		AppSecretKey:   key,
		AppSecretValue: kongKey.ID,
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
			// 启用密钥，调用 Kong Admin API 生成新的 secretValue
			kongKey, err := s.createKongKey(appKey.AppSecretKey)
			if err != nil {
				return err
			}
			appKey.AppSecretValue = kongKey.ID
		} else {
			// 禁用密钥，删除 Kong 中的记录
			err := s.deleteKongKey(appKey.AppSecretValue)
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

// 构建 Kong Admin API URL
func (s *AppKeyService) buildKongAdminURL(suffix string) string {
	return fmt.Sprintf("http://%s:%s/%s%s", s.host, s.port, s.path, suffix)
}

// 创建 Kong 密钥
func (s *AppKeyService) createKongKey(key string) (*KongSecretKeyVo, error) {
	url := s.buildKongAdminURL("")
	body := map[string]string{
		"key": key,
	}
	logx.Infof("buildKongAdminURL :%v", url)
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		logx.Errorf("Failed to create key in Kong Admin: %v", err)
		return nil, errors.NewCodeError(500, "创建密钥失败")
	}

	if resp.StatusCode()/100 != 2 {
		logx.Errorf("Failed to create key in Kong Admin: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, errors.NewCodeError(500, "创建密钥失败")
	}

	var kongKey KongSecretKeyVo
	err = json.Unmarshal(resp.Body(), &kongKey)
	if err != nil {
		logx.Errorf("Failed to parse Kong Admin response: %v, body: %s", err, resp.Body())
		return nil, errors.NewCodeError(500, "解析响应失败")
	}

	return &kongKey, nil
}

// 删除 Kong 密钥
func (s *AppKeyService) deleteKongKey(keyId string) error {
	url := s.buildKongAdminURL(fmt.Sprintf("/%s", keyId))

	resp, err := resty.New().R().
		Delete(url)

	if err != nil {
		logx.Errorf("Failed to delete key from Kong Admin: %v", err)
		return errors.NewCodeError(500, "删除密钥失败")
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 404 {
		logx.Errorf("Failed to delete key from Kong Admin: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return errors.NewCodeError(500, "删除密钥失败")
	}

	return nil
}

// 获取 Kong 密钥
func (s *AppKeyService) getKongKey(keyId string) (*KongSecretKeyVo, error) {
	url := s.buildKongAdminURL(fmt.Sprintf("/%s", keyId))

	resp, err := resty.New().R().
		Get(url)

	if err != nil {
		logx.Errorf("Failed to get key from Kong Admin: %v", err)
		return nil, errors.NewCodeError(500, "获取密钥失败")
	}

	if resp.StatusCode() == 404 {
		return nil, nil
	}

	if resp.StatusCode()/100 != 2 {
		logx.Errorf("Failed to get key from Kong Admin: status code %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, errors.NewCodeError(500, "获取密钥失败")
	}

	var kongKey KongSecretKeyVo
	err = json.Unmarshal(resp.Body(), &kongKey)
	if err != nil {
		logx.Errorf("Failed to parse Kong Admin response: %v, body: %s", err, resp.Body())
		return nil, errors.NewCodeError(500, "解析响应失败")
	}

	return &kongKey, nil
}
