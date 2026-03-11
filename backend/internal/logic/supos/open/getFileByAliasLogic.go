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

type GetFileByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 别名查询文件详情
func NewGetFileByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByAliasLogic {
	return &GetFileByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileByAliasLogic) GetFileByAlias(req *types.AliasPathReq) (resp *types.ResultVO, err error) {
	// 获取 uns 查询服务
	unsQueryService := spring.GetBean[*uns.UnsQueryService]()

	// 根据别名获取文件详情
	instanceDetail, err := unsQueryService.GetInstanceDetail(l.ctx, &types.InstanceDetailReq{}, req.Alias, "")
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

	// 转换 Fields 到 Definition
	convertedDetail := ConvertInstanceDetailToDefinition(instanceDetail.Data)

	// 返回成功结果
	return &types.ResultVO{
		Code: 200,
		Msg:  "成功",
		Data: convertedDetail,
	}, nil
}
