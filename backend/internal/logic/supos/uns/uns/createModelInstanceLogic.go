// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModelInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文件夹和文件
func NewCreateModelInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelInstanceLogic {
	return &CreateModelInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelInstanceLogic) CreateModelInstance(req *types.CreateTopicDto) (resp *types.CreateUnsResp, err error) {
	if alias := req.Alias; alias != "" {
		var mapper relationDB.UnsNamespaceRepo
		po, _ := mapper.GetByAlias(relationDB.GetDb(l.ctx), alias)
		if po != nil {
			var result types.CreateUnsResp
			result.Code = 400
			result.Msg = I18nUtils.GetMessageWithCtx(l.ctx, "uns.alias.has.exist")
			return &result, nil
		}
	}
	resp = spring.GetBean[*service.UnsAddService]().CreateModelInstance(l.ctx, req)
	return
}
