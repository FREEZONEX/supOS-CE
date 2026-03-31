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

type BatchQueryFileByPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据文件路径批量查询文件实时值
func NewBatchQueryFileByPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileByPathLogic {
	return &BatchQueryFileByPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileByPathLogic) BatchQueryFileByPath(req *types.StringArrayRequest) (resp *types.ResultVO, err error) {
	unsQueryService := spring.GetBean[*service.UnsQueryService]()

	resultMap := make(map[string]interface{})
	for _, path := range req.Items {
		// 获取最后一条消息
		msg, err := unsQueryService.GetLastMsgByPath(path)
		if err != nil {
			l.Errorf("查询文件 %s 失败: %v", path, err)
			resultMap[path] = nil
			continue
		}
		resultMap[path] = msg
	}

	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
		Data: resultMap,
	}, nil
}
