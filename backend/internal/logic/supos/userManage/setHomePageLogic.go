package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type SetHomePageLogic struct {
	baseUserManageLogic
}

// Set home page
func NewSetHomePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetHomePageLogic {
	return &SetHomePageLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *SetHomePageLogic) SetHomePage(req *types.HomePageReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.HomePage) == "" {
		return nil, errors.Parameter.WithMsg("homePage is required")
	}

	user := l.currentUser()
	if user == nil {
		return nil, errors.NotLogin
	}

	updateReq := &types.UserUpdateReq{
		UserID:   user.Sub,
		HomePage: strings.TrimSpace(req.HomePage),
	}
	result, err := NewUpdateUserLogic(l.ctx, l.svcCtx).UpdateUser(updateReq)
	if err != nil {
		return nil, err
	}

	user.HomePage = strings.TrimSpace(req.HomePage)
	l.updateUserCache(user)
	return result, nil
}
