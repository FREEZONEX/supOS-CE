package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type SetEmailLogic struct {
	baseUserManageLogic
}

// Update email address
func NewSetEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetEmailLogic {
	return &SetEmailLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *SetEmailLogic) SetEmail(req *types.EmailUpdateReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.Email) == "" {
		return nil, errors.Parameter.WithMsg("email is required")
	}

	user := l.currentUser()
	if user == nil {
		return nil, errors.NotLogin
	}

	updateReq := &types.UserUpdateReq{
		UserID: user.Sub,
		Email:  strings.TrimSpace(req.Email),
	}
	result, err := NewUpdateUserLogic(l.ctx, l.svcCtx).UpdateUser(updateReq)
	if err != nil {
		return nil, err
	}

	user.Email = strings.TrimSpace(req.Email)
	l.updateUserCache(user)
	return result, nil
}
