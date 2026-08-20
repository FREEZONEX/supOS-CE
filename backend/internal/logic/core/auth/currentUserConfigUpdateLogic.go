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

type CurrentUserConfigUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserConfigUpdateLogic {
	return &CurrentUserConfigUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserConfigUpdateLogic) CurrentUserConfigUpdate(req *types.CurrentUserConfigReq) (resp *types.Envelope, err error) {
	subject, ok := contextx.SubjectFrom(l.ctx)
	if !ok {
		return nil, respx.NewHTTPError(401, "unauthorized")
	}
	var homePage *string
	if req.HomePage != "" {
		resolvedHomePage, err := l.svcCtx.App.IAM.ResolveUserHomePage(l.ctx, subject.UserID, req.HomePage)
		if err != nil {
			return nil, err
		}
		homePage = &resolvedHomePage
	}
	var mainLanguage *string
	if req.MainLanguage != "" {
		mainLanguage = &req.MainLanguage
	}
	config, err := repo.NewUserConfigRepo(l.ctx).UpdateUserConfig(l.ctx, subject.UserID, homePage, mainLanguage)
	if err != nil {
		return nil, err
	}
	resolvedHomePage, err := l.svcCtx.App.IAM.ResolveUserHomePage(l.ctx, subject.UserID, config.HomePage)
	if err != nil {
		return nil, err
	}
	return respx.Envelope(types.CurrentUserConfigResp{
		UserId:       subject.UserID,
		HomePage:     resolvedHomePage,
		MainLanguage: config.MainLanguage,
	}), nil
}
