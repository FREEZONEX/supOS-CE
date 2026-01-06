// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"
	"encoding/json"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取最新消息
func NewGetLastMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMsgLogic {
	return &GetLastMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLastMsgLogic) GetLastMsg(req *types.GetLastMsgReq) (resp *types.GetLastMsgResp, err error) {
	sv := spring.GetBean[*service.UnsQueryService]()
	var bs []byte
	resp = &types.GetLastMsgResp{}
	resp.Code = 200
	resp.Msg = "OK"
	if id, alias := req.Id, req.Alias; id > 0 {
		bs, err = sv.GetLastMsg(id)
		resp.Data = json2Map(bs)
	} else if alias != "" {
		bs, err = sv.GetLastMsgByAlias(alias)
		resp.Data = json2Map(bs)
	} else if paths := req.Paths; len(paths) > 0 {
		combineMap := make(map[string]any)
		resp.Data = combineMap
		for _, path := range paths {
			bs, err = sv.GetLastMsgByPath(path)
			vm := json2Map(bs)
			if vm != nil {
				combineMap[path] = vm
			}
		}
	}
	return
}
func json2Map(bs []byte) (rs map[string]interface{}) {
	if len(bs) == 0 {
		return nil
	}
	er := json.Unmarshal(bs, &rs)
	if er != nil {
		return nil
	}
	return rs
}
