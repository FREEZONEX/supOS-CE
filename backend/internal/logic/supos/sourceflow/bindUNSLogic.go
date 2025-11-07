// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"

	unsservice "backend/internal/logic/supos/uns/uns/service"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type BindUNSLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	unsQueryService *unsservice.UnsQueryService
}

// bind a source flow with UNS alias
func NewBindUNSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUNSLogic {

	unsQueryService := &unsservice.UnsQueryService{}

	return &BindUNSLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		unsQueryService: unsQueryService,
	}
}

func (l *BindUNSLogic) BindUNS(req *types.SourceFlowBindUnsReq) error {
	if req.FlowID == 0 || req.UnsAlias == "" {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	db := stores.GetCommonConn(l.ctx)
	// 检查 UNS 是否存在
	var uns relationDB.UnsNamespace
	err := db.Where("alias", req.UnsAlias).First(&uns).Error
	if err != nil {
		return errors.Parameter.WithMsg("uns.file.not.exist")
	}
	// 删除旧的
	db.Where("alias", req.UnsAlias).Delete(&relationDB.NoderedFlowNode{})
	// 新增新的
	db.Create(&relationDB.NoderedFlowNode{
		Alias:    req.UnsAlias,
		ParentID: req.FlowID,
	})
	return nil
}
