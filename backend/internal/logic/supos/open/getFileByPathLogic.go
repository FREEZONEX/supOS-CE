// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	uns "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileByPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 路径查询文件详情
func NewGetFileByPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByPathLogic {
	return &GetFileByPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileByPathLogic) GetFileByPath(req *types.GetByPathReq) (resp *types.ResultVO, err error) {
	// 获取 uns 查询服务
	unsQueryService := spring.GetBean[*uns.UnsQueryService]()

	// 根据路径获取文件详情
	instanceDetail, err := unsQueryService.GetInstanceDetail(l.ctx, &types.InstanceDetailReq{}, "", req.Path)
	if err != nil {
		l.Errorf("获取文件详情失败: %v", err)
		return &types.ResultVO{
			Code: 500,
			Msg:  "获取文件详情失败",
			Data: nil,
		}, err
	}

	// 检查是否找到文件
	if instanceDetail.Data == nil {
		return &types.ResultVO{
			Code: 404,
			Msg:  "文件不存在",
			Data: nil,
		}, nil
	}

	// 返回成功结果
	return &types.ResultVO{
		Code: 200,
		Msg:  "成功",
		Data: instanceDetail.Data,
	}, nil
}
