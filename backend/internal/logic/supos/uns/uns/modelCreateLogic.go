package uns

/*
import (
	"context"

	"backend/internal/common/constants"
	"backend/internal/logic/supos/uns/common"
	"backend/internal/svc"
	"backend/internal/types"

	shareerrors "gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ModelCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文件夹和文件
func NewModelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelCreateLogic {
	return &ModelCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModelCreateLogic) ModelCreate(req *types.UnsCreateTopicDTO) (*types.WithID, error) {
	if req == nil {
		return nil, shareerrors.Parameter.WithMsg("请求参数为空")
	}

	var templateName string
	var err error
	if req.PathType == int64(constants.PathTypeDir) && req.CreateTemplate {
		template, err := common.CreateTemplate(l.ctx, l.svcCtx, req)
		if err != nil {
			return nil, err
		}
		if template != nil {
			req.ModelID = template.ID
			templateName = template.Name
		}
	}

	results, err := common.CreateUNSBatch(l.ctx, l.svcCtx, []*types.UnsCreateTopicDTO{req})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, shareerrors.System.WithMsg("创建UNS失败")
	}

	ns := results[0]
	if templateName != "" {
		l.Logger.Infof("目录[%s]已同步生成模板[%s]", ns.Name, templateName)
	}

	return &types.WithID{ID: ns.ID}, nil
}
*/
