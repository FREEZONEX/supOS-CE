package resource

import (
	"context"
	"strconv"
	"strings"
	"time"

	I18nUtils "backend/internal/common/I18nUtils"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List resources
func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.ResourceQuery) (resp []types.ResourceVO, err error) {
	db := stores.GetCommonConn(l.ctx)
	query := db.WithContext(l.ctx).Model(&relationDB.SuposResource{})
	if req != nil {
		if req.Type != 0 {
			query = query.Where("type = ?", req.Type)
		}
		if req.ParentID != 0 {
			query = query.Where("parent_id = ?", req.ParentID)
		}
	}
	var records []relationDB.SuposResource
	if err = query.Order("sort ASC, id ASC").Find(&records).Error; err != nil {
		l.Errorf("failed to query resources: %v", err)
		return nil, errors.Database.WithMsg("resource.query.failed").AddDetail(err)
	}
	resp = make([]types.ResourceVO, 0, len(records))
	for _, item := range records {
		resp = append(resp, toResourceVO(l.ctx, &item))
	}
	return resp, nil
}

func toResourceVO(ctx context.Context, res *relationDB.SuposResource) types.ResourceVO {
	if res == nil {
		return types.ResourceVO{}
	}
	var parentID string
	if res.ParentID != nil && *res.ParentID > 0 {
		parentID = strconv.FormatInt(*res.ParentID, 10)
	}
	showName := I18nUtils.GetMessageWithCtx(ctx, stringValue(res.NameCode))
	if showName == "" {
		showName = stringValue(res.NameCode)
	}
	showName = transI18nSTR(showName)
	showDesc := I18nUtils.GetMessageWithCtx(ctx, stringValue(res.DescriptionCode))
	if showDesc == "" {
		showDesc = stringValue(res.DescriptionCode)
	}
	showDesc = transI18nSTR(showDesc)
	return types.ResourceVO{
		ID:              strconv.FormatInt(res.ID, 10),
		ParentID:        parentID,
		Type:            res.Type,
		Code:            res.Code,
		NameCode:        stringValue(res.NameCode),
		ShowName:        showName,
		RouteSource:     derefInt(res.RouteSource),
		URL:             stringValue(res.URL),
		URLType:         derefInt(res.URLType),
		OpenType:        derefInt(res.OpenType),
		Icon:            stringValue(res.Icon),
		DescriptionCode: stringValue(res.DescriptionCode),
		ShowDescription: showDesc,
		Sort:            derefInt(res.Sort),
		EditEnable:      boolValue(res.EditEnable),
		HomeEnable:      boolValue(res.HomeEnable),
		Fixed:           boolValue(res.Fixed),
		Enable:          boolValue(res.Enable),
		UpdateAt:        formatTime(res.UpdateAt),
		CreateAt:        formatTime(res.CreateAt),
	}
}

func transI18nSTR(str string) string {
	indent := "%!(EXTRA []interface {}=[])"
	if strings.Contains(str, indent) {
		return strings.ReplaceAll(str, indent, "")
	}
	return str
}

func derefInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func stringPtr(val string) *string {
	v := strings.TrimSpace(val)
	if v == "" {
		return nil
	}
	return &v
}

func intPtr(val int) *int {
	return &val
}

func int64Ptr(val int64) *int64 {
	return &val
}

func optionalParentID(parentID int64) *int64 {
	if parentID <= 0 {
		return nil
	}
	return &parentID
}

func boolPtr(val bool) *bool {
	return &val
}

func stringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func stringValueForUpdate(ptr *string) any {
	if ptr == nil {
		return nil
	}
	return *ptr
}

func boolValue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}
