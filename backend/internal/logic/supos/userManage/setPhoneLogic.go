package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type SetPhoneLogic struct {
	baseUserManageLogic
}

// Update phone number
func NewSetPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPhoneLogic {
	return &SetPhoneLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *SetPhoneLogic) SetPhone(req *types.PhoneUpdateReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.Phone) == "" {
		return nil, errors.Parameter.WithMsg("phone is required")
	}

	user := l.currentUser()
	if user == nil {
		return nil, errors.NotLogin
	}

	updateReq := &types.UserUpdateReq{
		UserID: user.Sub,
		Phone:  strings.TrimSpace(req.Phone),
	}
	result, err := NewUpdateUserLogic(l.ctx, l.svcCtx).UpdateUser(updateReq)
	if err != nil {
		return nil, err
	}

	user.Phone = strings.TrimSpace(req.Phone)
	l.updateUserCache(user)
	return result, nil
}
