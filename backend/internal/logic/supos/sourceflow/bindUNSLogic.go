// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/share/spring"
	"context"
	"time"

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

	tx := db.Begin()
	defer tx.Rollback()

	// 删除旧的
	err = tx.Where("alias", req.UnsAlias).Delete(&relationDB.NoderedFlowNode{}).Error
	if err != nil {
		return errors.Parameter.WithMsg("uns.file.not.exist")
	}
	// 新增新的
	tx.Create(&relationDB.NoderedFlowNode{
		Alias:    req.UnsAlias,
		ParentID: req.FlowID,
	})
	var unsMapper relationDB.UnsNamespaceRepo
	// 更新 UNS flag
	flag := int32(constants.UnsFlagWithFlow)
	if flags := uns.WithFlags; flags != nil {
		flag = *flags | int32(constants.UnsFlagWithFlow)
	}
	uns.WithFlags = &flag
	uns.UpdateAt = time.Now()
	err = unsMapper.Update(tx, &uns)
	if err != nil {
		return errors.Parameter.WithMsg(err.Error())
	}
	tx.Commit()
	err = spring.PublishEvent(&event.UpdateInstanceEvent{ApplicationEvent: event.ApplicationEvent{Context: l.ctx},
		Topics: []*types.CreateTopicDto{UnsConverter.Po2Dto(&uns)}})
	return nil
}
