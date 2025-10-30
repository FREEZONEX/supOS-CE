// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/logic"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存菜单
func NewSaveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMenuLogic {
	return &SaveMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveMenuLogic) SaveMenu(menuDto *dto.MenuDto) (resp *types.ResultVO, err error) {
	menuLogic := logic.NewMenuLogic(l.svcCtx.Config.Kong.Host, l.svcCtx.Config.Kong.Port)
	err = menuLogic.CreateRoute(menuDto, false, false)
	if err != nil {
		return nil, err
	}
	resp = &types.ResultVO{
		Code: 0,
		Msg:  "success",
	}
	return
}
