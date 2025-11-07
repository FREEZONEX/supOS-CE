package person

import (
	"context"
	"strings"
	"time"

	"backend/internal/common/utils/langutil"
	"backend/internal/common/vo"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取个人配置
func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config(req *types.GetPersonConfigReq) (*types.GetPersonConfigResp, error) {
	if req == nil {
		req = &types.GetPersonConfigReq{}
	}

	user := currentUser(l.ctx)
	if user == nil || strings.TrimSpace(user.Sub) == "" {
		return nil, errors.NotLogin
	}

	targetUserID := strings.TrimSpace(req.UserID)
	if targetUserID == "" {
		targetUserID = strings.TrimSpace(user.Sub)
	} else if !strings.EqualFold(targetUserID, user.Sub) {
		return nil, errors.Permissions.WithMsg("common.noPermissionMessage")
	}

	repo := relationDB.NewUnsPersonConfigRepo(l.ctx)
	cfg, err := repo.FindOneByFilter(l.ctx, relationDB.UnsPersonConfigFilter{UserID: targetUserID})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		cfg = nil
	}
	resp := &types.GetPersonConfigResp{
		UserID:       targetUserID,
		MainLanguage: preferredLanguage(cfg, user),
	}
	if cfg != nil {
		resp.CreateAt = unixMilliOrZero(cfg.CreateAt)
		resp.UpdateAt = unixMilliOrZero(cfg.UpdateAt)
	}

	return resp, nil
}

func preferredLanguage(cfg *relationDB.UnsPersonConfig, user *vo.UserInfoVo) string {
	candidates := []string{}
	if cfg != nil {
		candidates = append(candidates, cfg.MainLanguage)
	}
	if user != nil {
		candidates = append(candidates, user.MainLanguage)
	}
	candidates = append(candidates, langutil.SystemLocale())

	for _, lang := range candidates {
		if val := strings.TrimSpace(lang); val != "" {
			return val
		}
	}
	return langutil.SystemLocale()
}

func unixMilliOrZero(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}
