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

type BatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Batch delete resources
func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLogic) BatchDelete(req *types.ResourceBatchDeleteReq) error {
	if req == nil || len(req.IDs) == 0 {
		return errors.Parameter.WithMsg("resource.batch.delete.empty")
	}
	db := stores.GetCommonConn(l.ctx)
	return db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", req.IDs).Delete(&relationDB.SuposResource{}).Error; err != nil {
			l.Errorf("failed to batch delete resources: %v", err)
			return errors.Database.WithMsg("resource.delete.failed").AddDetail(err)
		}
		if err := tx.Where("parent_id IN ?", req.IDs).Delete(&relationDB.SuposResource{}).Error; err != nil {
			l.Errorf("failed to delete child resources: %v", err)
			return errors.Database.WithMsg("resource.delete.failed").AddDetail(err)
		}
		return nil
	})
}
