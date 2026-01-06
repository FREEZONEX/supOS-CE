package resource

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete resource by id
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.ResourceIDReq) error {
	if req == nil || req.ID == 0 {
		return errors.Parameter.WithMsg("resource.id.not.found")
	}
	db := stores.GetCommonConn(l.ctx)
	return db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", req.ID).Delete(&relationDB.SuposResource{}).Error; err != nil {
			l.Errorf("failed to delete resource %d: %v", req.ID, err)
			return errors.Database.WithMsg("resource.delete.failed").AddDetail(err)
		}
		if err := tx.Where("parent_id = ?", req.ID).Delete(&relationDB.SuposResource{}).Error; err != nil {
			l.Errorf("failed to delete child resources for %d: %v", req.ID, err)
			return errors.Database.WithMsg("resource.delete.failed").AddDetail(err)
		}
		return nil
	})
}
