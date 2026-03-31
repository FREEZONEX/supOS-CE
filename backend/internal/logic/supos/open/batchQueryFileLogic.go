// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件实时值
func NewBatchQueryFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileLogic {
	return &BatchQueryFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileLogic) BatchQueryFile(req *types.StringArrayRequest) (resp *types.ResultVO, err error) {
	unsQueryService := spring.GetBean[*service.UnsQueryService]()

	resultMap := make(map[string]interface{})
	for _, alias := range req.Items {
		// 获取最后一条消息
		msg, err := unsQueryService.GetLastMsgByAlias(alias)
		if err != nil {
			l.Errorf("查询文件 %s 失败: %v", alias, err)
			resultMap[alias] = nil
			continue
		}
		resultMap[alias] = msg
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
		Data: resultMap,
	}, nil
}
