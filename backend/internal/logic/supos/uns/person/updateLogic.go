package person

import (
	"context"
	"strings"
	"time"

	cache "backend/internal/common/cache"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置个人配置
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdatePersonConfigReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	lang := strings.TrimSpace(req.MainLanguage)
	if lang == "" {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}

	user := currentUser(l.ctx)
	if user == nil || strings.TrimSpace(user.Sub) == "" {
		return errors.NotLogin
	}

	targetUserID := strings.TrimSpace(req.UserID)
	if targetUserID == "" {
		targetUserID = user.Sub
	} else if !strings.EqualFold(targetUserID, user.Sub) {
		return errors.Permissions.WithMsg("common.noPermissionMessage")
	}

	repo := relationDB.NewUnsPersonConfigRepo(l.ctx)
	cfg, err := repo.FindOneByFilter(l.ctx, relationDB.UnsPersonConfigFilter{UserID: targetUserID})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return err
	}

	now := time.Now()
	db := stores.GetCommonConn(l.ctx)
	if db == nil {
		return errors.System.WithMsg("common database connection not initialized")
	}

	if err := db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if cfg == nil {
			cfg = &relationDB.UnsPersonConfig{
				UserID:       targetUserID,
				MainLanguage: lang,
				CreateAt:     now,
				UpdateAt:     now,
			}
			if err := tx.Create(cfg).Error; err != nil {
				return stores.ErrFmt(err)
			}
		} else {
			cfg.MainLanguage = lang
			cfg.UpdateAt = now
			if err := tx.Model(&relationDB.UnsPersonConfig{}).Where("id = ?", cfg.ID).Save(cfg).Error; err != nil {
				return stores.ErrFmt(err)
			}
		}

		if err := tx.Model(&relationDB.IamUser{}).Where("id = ?", targetUserID).Updates(map[string]any{
			"main_language": lang,
			"updated_at":    now,
		}).Error; err != nil {
			return stores.ErrFmt(err)
		}

		return nil
	}); err != nil {
		return err
	}

	user.MainLanguage = lang
	if cache.UserInfoCache != nil {
		cache.UserInfoCache.Set(user.Sub, user)
	}

	return nil
}
