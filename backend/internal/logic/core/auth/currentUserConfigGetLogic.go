package auth

import (
	"context"

	"backend/internal/contextx"
	respx "backend/internal/httpx"
	"backend/internal/repo"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserConfigGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserConfigGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserConfigGetLogic {
	return &CurrentUserConfigGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserConfigGetLogic) CurrentUserConfigGet() (resp *types.Envelope, err error) {
	subject, ok := contextx.SubjectFrom(l.ctx)
	if !ok {
		return nil, respx.NewHTTPError(401, "unauthorized")
	}
	config, err := repo.NewUserConfigRepo(l.ctx).GetUserConfig(l.ctx, subject.UserID)
	if err != nil {
		return nil, err
	}
	homePage, err := l.svcCtx.App.IAM.ResolveUserHomePage(l.ctx, subject.UserID, config.HomePage)
	if err != nil {
		return nil, err
	}
	return respx.Envelope(types.CurrentUserConfigResp{
		UserId:       subject.UserID,
		HomePage:     homePage,
		MainLanguage: config.MainLanguage,
	}), nil
}
