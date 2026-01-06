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
	"github.com/zeromicro/go-zero/core/logx"
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
	if cfg == nil {
		cfg = &relationDB.UnsPersonConfig{
			UserID:       targetUserID,
			MainLanguage: lang,
			CreateAt:     now,
			UpdateAt:     now,
		}
		if err := repo.Insert(l.ctx, cfg); err != nil {
			return err
		}
	} else {
		cfg.MainLanguage = lang
		cfg.UpdateAt = now
		if err := repo.Update(l.ctx, cfg); err != nil {
			return err
		}
	}

	user.MainLanguage = lang
	if cache.UserInfoCache != nil {
		cache.UserInfoCache.Set(user.Sub, user)
	}

	return nil
}
