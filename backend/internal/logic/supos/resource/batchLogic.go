package resource

import (
	"context"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
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

func (l *BatchLogic) Batch(req *[]types.BatchUpdateResource) error {
	if req == nil || len(*req) == 0 {
		return errors.Parameter.WithMsg("resource.batch.empty")
	}
	db := stores.GetCommonConn(l.ctx)
	return db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range *req {
			if item.ID == 0 {
				return errors.Parameter.WithMsg("resource.id.not.found")
			}
			updates := make(map[string]any)

			if item.Type != nil {
				updates["type"] = *item.Type
			}
			if item.Source != nil {
				updates["source"] = stringValueForUpdate(item.Source)
			}
			if item.Code != nil {
				updates["code"] = stringValueForUpdate(item.Code)
			}
			if item.NameCode != nil {
				updates["name_code"] = stringValueForUpdate(item.NameCode)
			}
			if item.Name != nil {
				updates["name"] = stringValueForUpdate(item.Name)
			}
			if item.RouteSource != nil {
				updates["route_source"] = *item.RouteSource
			}
			if item.URL != nil {
				updates["url"] = stringValueForUpdate(item.URL)
			}
			if item.URLType != nil {
				updates["url_type"] = *item.URLType
			}
			if item.OpenType != nil {
				updates["open_type"] = *item.OpenType
			}
			if item.Icon != nil {
				updates["icon"] = stringValueForUpdate(item.Icon)
			}
			if item.DescriptionCode != nil {
				updates["description_code"] = stringValueForUpdate(item.DescriptionCode)
			} else if item.Description != nil {
				updates["description_code"] = stringValueForUpdate(item.Description)
			}
			if item.Sort != nil {
				updates["sort"] = *item.Sort
			}
			if item.EditEnable != nil {
				updates["edit_enable"] = *item.EditEnable
			}
			if item.HomeEnable != nil {
				updates["home_enable"] = *item.HomeEnable
			}
			if item.Fixed != nil {
				updates["fixed"] = *item.Fixed
			}
			if item.Enable != nil {
				updates["enable"] = *item.Enable
			}
			if item.ParentID != nil {
				if *item.ParentID <= 0 {
					updates["parent_id"] = nil
				} else {
					updates["parent_id"] = *item.ParentID
				}
			}

			if len(updates) == 0 {
				continue
			}
			updates["update_at"] = time.Now()

			if err := tx.Model(&relationDB.SuposResource{}).Where("id = ?", item.ID).Updates(updates).Error; err != nil {
				l.Errorf("batch update resource %d failed: %v", item.ID, err)
				return errors.Database.WithMsg("resource.batch.update.failed").AddDetail(err)
			}
		}
		return nil
	})
}
