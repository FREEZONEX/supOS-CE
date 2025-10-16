package label

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"strings"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.UnsLabelCreateReq) (resp *types.WithID, err error) {
	// 参数校验
	if strings.TrimSpace(req.LabelName) == "" {
		return nil, errors.Parameter.WithMsg("labelName不能为空")
	}

	// 写入数据库
	db := relationDB.NewUnsLabelRepo(stores.GetCommonConn(l.ctx))
	data := &relationDB.UnsLabel{
		LabelName: req.LabelName,
	}
	if err = db.Insert(l.ctx, data); err != nil {
		return nil, err
	}

	// 返回新ID
	return &types.WithID{ID: data.ID}, nil
}
