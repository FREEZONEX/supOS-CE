// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/uns/service"
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryFileHistoryValueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件历史值
func NewBatchQueryFileHistoryValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileHistoryValueLogic {
	return &BatchQueryFileHistoryValueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileHistoryValueLogic) BatchQueryFileHistoryValue(req *types.HistoryValueRequest) (resp *types.UnsHistoryQueryResult, err error) {
	resp = &types.UnsHistoryQueryResult{
		BaseResult: types.BaseResult{
			Code: 200,
			Msg:  "ok",
		},
	}
	if req == nil || !hasHistoryAliases(req.AliasList) {
		resp.Code = 400
		resp.Msg = I18nUtils.GetMessageWithCtx(l.ctx, "uns.file.history.query.alias.required")
		return resp, nil
	}
	if req.TimeStart <= 0 || req.TimeEnd <= 0 || req.TimeEnd < req.TimeStart {
		resp.Code = 400
		resp.Msg = I18nUtils.GetMessageWithCtx(l.ctx, "uns.file.history.query.time.invalid")
		return resp, nil
	}

	data, err := service.NewHistoryQueryService().Query(l.ctx, req)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	resp.Data = data
	if data != nil && (len(data.NotExists) > 0 || len(data.ErrorFields) > 0) {
		resp.Code = 206
	}
	return resp, nil
}

func hasHistoryAliases(aliasList []string) bool {
	for _, alias := range aliasList {
		if strings.TrimSpace(alias) != "" {
			return true
		}
	}
	return false
}
