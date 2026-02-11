// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/eventflow"
	"backend/internal/logic/supos/sourceflow"
	dashboardLogic "backend/internal/logic/supos/uns/dashboard"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"strconv"

	"gitee.com/unitedrhino/share/stores"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID删除组
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.GroupIDReq) (*types.JsonResult, error) {
	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	group, err := groupMapper.SelectById(db, req.ID)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  I18nUtils.GetMessage("group.query.failed") + ": " + err.Error(),
		}, nil
	}
	if group == nil {
		return &types.JsonResult{
			Code: 400,
			Msg:  I18nUtils.GetMessage("group.notfound"),
		}, nil
	}

	// 查询分组下的所有业务数据
	sourceFlows, eventFlows, dashboards, err := groupMapper.CleanBizGroupById(db, req.ID, group.Type)
	if err != nil {
		l.Errorf("查询分组关联的业务数据失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  I18nUtils.GetMessage("group.cleanbiz.query.failed") + ": " + err.Error(),
		}, nil
	}

	// 根据分组类型删除对应的业务数据
	switch *group.Type {
	case 1: // groupType 1 对应 supos_node_flows 表 (source flow)
		// 删除 source flow
		for _, flow := range sourceFlows {
			if flow.FlowID == "" {
				continue
			}
			req := &types.SourceFlowDeleteReq{
				ID: strconv.FormatInt(flow.ID, 10),
			}
			err = sourceflow.NewDeleteSourceFlowLogic(l.ctx, l.svcCtx).DeleteSourceFlow(req)
			if err != nil {
				l.Errorf("删除 source flow 失败: %v", err)
			}
		}

	case 2: // groupType 2 对应 supos_node_flows 表 (event flow)
		// 删除 event flow
		for _, flow := range eventFlows {
			req := &types.EventFlowDeleteReq{
				ID: strconv.FormatInt(flow.ID, 10),
			}
			err = eventflow.NewDeleteEventFlowLogic(l.ctx, l.svcCtx).DeleteEventFlow(req)
			if err != nil {
				l.Errorf("删除 event flow 失败: %v", err)
			}
		}

	case 3: // groupType 3 对应 uns_dashboard 表 (dashboard)
		// 删除 dashboard
		for _, d := range dashboards {
			_, err = dashboardLogic.NewDeleteLogic(l.ctx, l.svcCtx).Delete(d.ID)
			if err != nil {
				l.Errorf("删除 dashboard 失败: %v", err)
			}
		}

	default:
		l.Infof("unknown group type: %d", *group.Type)
	}

	// 执行删除分组操作
	if err = groupMapper.DeleteById(db, req.ID); err != nil {
		l.Errorf("删除组失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  I18nUtils.GetMessage("group.delete.failed") + ": " + err.Error(),
		}, nil
	}

	return &types.JsonResult{
		Code: 200,
		Msg:  I18nUtils.GetMessage("group.delete.success"),
		Data: nil,
	}, nil
}
