package resource

import (
	"context"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/i18ns"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type BatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Batch update resources
func NewBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchLogic {
	return &BatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchLogic) Batch(req *types.BatchUpdateResourceReq) error {
	if req == nil || len(req.Items) == 0 {
		return errors.Parameter.WithMsg(i18ns.LocalizeMsg("resource.batch.empty"))
	}
	db := stores.GetCommonConn(l.ctx)
	return db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			if item.ID == 0 {
				return errors.Parameter.WithMsg(i18ns.LocalizeMsg("resource.id.not.found"))
			}
			source := stringPtr(item.Source)
			nameCode := stringPtr(item.Name)
			descCode := stringPtr(item.Description)
			url := stringPtr(item.URL)
			icon := stringPtr(item.Icon)

			updates := map[string]any{
				"type":             item.Type,
				"source":           stringValueForUpdate(source),
				"code":             item.Code,
				"name_code":        stringValueForUpdate(nameCode),
				"route_source":     item.RouteSource,
				"url":              stringValueForUpdate(url),
				"url_type":         item.URLType,
				"open_type":        item.OpenType,
				"icon":             stringValueForUpdate(icon),
				"description_code": stringValueForUpdate(descCode),
				"sort":             item.Sort,
				"edit_enable":      item.EditEnable,
				"home_enable":      item.HomeEnable,
				"fixed":            item.Fixed,
				"enable":           item.Enable,
				"update_at":        time.Now(),
			}
			if item.ParentID <= 0 {
				updates["parent_id"] = nil
			} else {
				updates["parent_id"] = item.ParentID
			}

			if err := tx.Model(&relationDB.SuposResource{}).Where("id = ?", item.ID).Updates(updates).Error; err != nil {
				l.Errorf("batch update resource %d failed: %v", item.ID, err)
				return errors.Database.WithMsg(i18ns.LocalizeMsg("resource.batch.update.failed")).AddDetail(err)
			}
		}
		return nil
	})
}
