// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"

	domainiam "backend/internal/domain/iam"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserProfileUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserProfileUpdateLogic {
	return &CurrentUserProfileUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserProfileUpdateLogic) CurrentUserProfileUpdate(req *types.CurrentUserProfileReq) (resp *types.Envelope, err error) {
	userID := logicx.UserID(l.ctx)
	if userID <= 0 {
		return nil, respx.NewHTTPError(401, "unauthorized")
	}
	data, err := l.svcCtx.App.IAM.UpdateUser(l.ctx, domainiam.UserUpdateCommand{
		ID:        userID,
		FirstName: req.FirstName,
		Email:     req.Email,
		Phone:     req.Phone,
		Enabled:   true,
		UserID:    userID,
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
