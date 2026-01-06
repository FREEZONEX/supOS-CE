package userManage

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type SetTipsEnableLogic struct {
	baseUserManageLogic
}

// Toggle tips flag
func NewSetTipsEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTipsEnableLogic {
	return &SetTipsEnableLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *SetTipsEnableLogic) SetTipsEnable(req *types.TipsEnableReq) (*types.OperationResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("tipsEnable is required")
	}

	user := l.currentUser()
	if user == nil {
		return nil, errors.NotLogin
	}

	value := req.TipsEnable
	updateReq := &types.UserUpdateReq{
		UserID:     user.Sub,
		TipsEnable: &value,
	}
	result, err := NewUpdateUserLogic(l.ctx, l.svcCtx).UpdateUser(updateReq)
	if err != nil {
		return nil, err
	}

	user.TipsEnable = value
	l.updateUserCache(user)
	return result, nil
}
