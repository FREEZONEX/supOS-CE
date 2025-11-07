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

type PostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create or update resource and its children
func NewPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostLogic {
	return &PostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostLogic) Post(req *types.SaveResourceReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("resource.request.empty")
	}

	db := stores.GetCommonConn(l.ctx)
	return db.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		_, err := l.saveResource(tx, req)
		return err
	})
}

func (l *PostLogic) saveResource(tx *gorm.DB, dto *types.SaveResourceReq) (int64, error) {
	now := time.Now()
	parentID := optionalParentID(dto.ParentID)
	nameCode := stringPtr(dto.Name)
	descCode := stringPtr(dto.Description)
	source := stringPtr(dto.Source)
	url := stringPtr(dto.URL)
	icon := stringPtr(dto.Icon)

	if dto.ID == 0 {
		var count int64
		if err := tx.Model(&relationDB.SuposResource{}).Where("code = ?", dto.Code).Count(&count).Error; err != nil {
			l.Errorf("failed to check resource code: %v", err)
			return 0, errors.Database.WithMsg("resource.save.failed").AddDetail(err)
		}
		if count > 0 {
			return 0, errors.Duplicate.WithMsg("resource.code.duplicate").AddDetail(dto.Code)
		}
		flag := true
		res := &relationDB.SuposResource{
			ParentID:        parentID,
			Type:            dto.Type,
			Source:          source,
			Code:            dto.Code,
			NameCode:        nameCode,
			RouteSource:     intPtr(dto.RouteSource),
			URL:             url,
			URLType:         intPtr(dto.URLType),
			OpenType:        intPtr(dto.OpenType),
			Icon:            icon,
			DescriptionCode: descCode,
			Sort:            intPtr(dto.Sort),
			EditEnable:      &flag,
			HomeEnable:      boolPtr(dto.HomeEnable),
			Fixed:           boolPtr(dto.Fixed),
			Enable:          &flag,
			CreateAt:        now,
			UpdateAt:        now,
		}
		if err := tx.Create(res).Error; err != nil {
			l.Errorf("failed to insert resource: %v", err)
			return 0, errors.Database.WithMsg("resource.save.failed").AddDetail(err)
		}
		id := res.ID
		for i := range dto.Children {
			dto.Children[i].ParentID = id
			if _, err := l.saveResource(tx, &dto.Children[i]); err != nil {
				return 0, err
			}
		}
		return id, nil
	}

	var existing relationDB.SuposResource
	if err := tx.Where("id = ?", dto.ID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.NotFind.WithMsg("resource.id.not.found").AddDetail(dto.ID)
		}
		l.Errorf("failed to load resource %d: %v", dto.ID, err)
		return 0, errors.Database.WithMsg("resource.save.failed").AddDetail(err)
	}

	updates := map[string]any{
		"type":             dto.Type,
		"source":           stringValueForUpdate(source),
		"code":             dto.Code,
		"name_code":        stringValueForUpdate(nameCode),
		"route_source":     dto.RouteSource,
		"url":              stringValueForUpdate(url),
		"url_type":         dto.URLType,
		"open_type":        dto.OpenType,
		"icon":             stringValueForUpdate(icon),
		"description_code": stringValueForUpdate(descCode),
		"sort":             dto.Sort,
		"home_enable":      dto.HomeEnable,
		"fixed":            dto.Fixed,
		"update_at":        now,
	}
	if parentID == nil {
		updates["parent_id"] = nil
	} else {
		updates["parent_id"] = *parentID
	}

	if err := tx.Model(&relationDB.SuposResource{}).Where("id = ?", dto.ID).Updates(updates).Error; err != nil {
		l.Errorf("failed to update resource %d: %v", dto.ID, err)
		return 0, errors.Database.WithMsg("resource.save.failed").AddDetail(err)
	}
	for i := range dto.Children {
		dto.Children[i].ParentID = dto.ID
		if _, err := l.saveResource(tx, &dto.Children[i]); err != nil {
			return 0, err
		}
	}
	return dto.ID, nil
}
