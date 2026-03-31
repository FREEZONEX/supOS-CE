// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	uns "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFolderByPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 路径查询文件夹详情
func NewGetFolderByPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFolderByPathLogic {
	return &GetFolderByPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFolderByPathLogic) GetFolderByPath(req *types.GetByPathReq) (resp *types.ResultVO, err error) {
	// 获取 uns 查询服务
	unsQueryService := spring.GetBean[*uns.UnsQueryService]()

	// 根据路径获取文件夹详情
	modelDefinition, err := unsQueryService.GetModelDefinition(l.ctx, &types.ModelDetailReq{}, "", req.Path)
	if err != nil {
		l.Errorf("获取文件夹详情失败: %v", err)
		return &types.ResultVO{
			Code: 500,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.operation.failed"),
			Data: nil,
		}, err
	}

	// 检查是否找到文件夹
	if modelDefinition.Data == nil {
		return &types.ResultVO{
			Code: 404,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.folder.not.found"),
			Data: nil,
		}, nil
	}

	// 转换 Fields 到 Definition
	convertedDetail := ConvertModelDetailToDefinition(modelDefinition.Data)

	// 返回成功结果
	return &types.ResultVO{
		Code: 200,
		Msg:  "ok",
		Data: convertedDetail,
	}, nil
}
