// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/logic"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/app/adapter"
	"backend/share/app/model"
	"backend/share/base"
	"context"

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
	menuLogic := logic.NewMenuLogic(l.svcCtx)
	err = menuLogic.CreateMenu(l.ctx, menuDto, true)
	if err != nil {
		return nil, err
	}
	serviceName, name := menuDto.ServiceName, menuDto.Name
	err = adapter.SaveMenuToDatabase(&model.MenuModel{
		Name:        name,
		Description: menuDto.Description,
		IndexUrl:    menuDto.BaseURL,
		OpenType:    menuDto.OpenType,
		IconUrl:     menuDto.Icon.Filename,
	}, base.SanYuan(len(serviceName) > 0, serviceName, name))
	if err != nil {
		resp = &types.ResultVO{
			Code: 500,
			Msg:  err.Error(),
		}
	} else {
		resp = &types.ResultVO{
			Code: 200,
			Msg:  "success",
		}
	}
	return
}
