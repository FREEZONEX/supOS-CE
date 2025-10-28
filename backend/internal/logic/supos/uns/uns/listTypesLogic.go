// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/share/base"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 枚举数据类型
func NewListTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTypesLogic {
	return &ListTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTypesLogic) ListTypes() (resp *types.ListTypesResult, err error) {
	return &types.ListTypesResult{
		BaseResult: types.BaseResult{Code: 200, Msg: "ok"},
		Data: base.Map[types.FieldType, string](types.FieldTypes(), func(e types.FieldType) string {
			return e.Name()
		}),
	}, nil
}
