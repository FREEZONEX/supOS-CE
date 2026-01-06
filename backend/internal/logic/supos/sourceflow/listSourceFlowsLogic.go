// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"backend/internal/repo/relationDB"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/common/constants"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ListSourceFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List source flows with optional fuzzy search
func NewListSourceFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSourceFlowsLogic {
	return &ListSourceFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSourceFlowsLogic) ListSourceFlows(req *types.SourceFlowListQuery) (*types.SourceFlowPageResult, error) {
	return l.ListFlowsWithType(req, constants.FlowTypeNODERED)
}

// ListFlowsWithType lists flows filtered by the provided flow type(template).
func (l *ListSourceFlowsLogic) ListFlowsWithType(req *types.SourceFlowListQuery, flowType string) (*types.SourceFlowPageResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request is nil")
	}
	pageNo := req.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	page := &stores.PageInfo{Page: pageNo, Size: pageSize}
	keyword := strings.TrimSpace(req.Keyword)

	template := strings.TrimSpace(flowType)
	if template == "" {
		template = constants.FlowTypeNODERED
	}
	userID := resolveUserID(l.ctx)

	db := stores.GetCommonConn(l.ctx)
	query := db.Table("supos_node_flows AS f").
		Select("f.id, f.flow_id, f.flow_name, f.flow_status, f.template, f.description, f.create_time,f.creator, f.update_time, COALESCE(t.mark, 0) AS mark, t.mark_time").
		Joins("LEFT JOIN supos_node_flow_top_recodes AS t ON f.id = t.id AND t.user_id = ?", userID).
		Where("f.template = ?", template)
	if keyword != "" {
		like := "%" + relationDB.EscapeLike(keyword) + "%"
		query = query.Where("(f.flow_name LIKE ? OR f.description LIKE ?)", like, like)
	}
	countQuery := query.Session(&gorm.Session{}).Select("f.id")
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	query = page.ToGorm(query)
	query = applyFlowOrdering(query, strings.TrimSpace(req.OrderCode), strings.TrimSpace(req.IsAsc))

	var rows []flowListItem
	if err := query.Find(&rows).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	items := make([]types.SourceFlowInfo, 0, len(rows))
	for _, row := range rows {
		mark := 0
		if row.Mark.Valid {
			mark = int(row.Mark.Int64)
		}
		items = append(items, types.SourceFlowInfo{
			ID:          strconv.FormatInt(row.ID, 10),
			FlowName:    row.FlowName,
			FlowID:      row.FlowID,
			Description: row.Description,
			FlowStatus:  row.FlowStatus,
			Template:    row.Template,
			Mark:        mark,
			Creator:     row.Creator,
			CreateTime:  int(row.CreateTime.UnixMilli()),
		})
	}
	return &types.SourceFlowPageResult{
		Code:     http.StatusOK,
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Data:     items,
	}, nil
}

type flowListItem struct {
	ID          int64         `gorm:"column:id"`
	FlowID      string        `gorm:"column:flow_id"`
	FlowName    string        `gorm:"column:flow_name"`
	FlowStatus  string        `gorm:"column:flow_status"`
	Template    string        `gorm:"column:template"`
	Description string        `gorm:"column:description"`
	CreateTime  time.Time     `gorm:"column:create_time"`
	UpdateTime  time.Time     `gorm:"column:update_time"`
	Mark        sql.NullInt64 `gorm:"column:mark"`
	MarkTime    *time.Time    `gorm:"column:mark_time"`
	Creator     string        `gorm:"column:creator"`
}

func applyFlowOrdering(db *gorm.DB, orderCode, isAsc string) *gorm.DB {
	if db == nil {
		return db
	}
	db = db.Order("COALESCE(t.mark, 0) DESC")
	column := normalizeOrderColumn(orderCode)
	if column == "" {
		return db.Order("t.mark_time DESC").Order("f.create_time DESC")
	}
	direction := "DESC"
	if strings.EqualFold(isAsc, "true") {
		direction = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", column, direction))
}

func normalizeOrderColumn(orderCode string) string {
	if orderCode == "" {
		return ""
	}
	snake := camelToSnake(orderCode)
	if snake == "" {
		return ""
	}
	switch snake {
	case "mark_time":
		return "t.mark_time"
	case "mark":
		return "COALESCE(t.mark, 0)"
	case "flow_name", "flow_status", "description", "create_time", "update_time", "template":
		return "f." + snake
	default:
		return ""
	}
}

var (
	camelRegex      = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snakeValidRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
)

func camelToSnake(input string) string {
	if input == "" {
		return ""
	}
	snake := camelRegex.ReplaceAllString(input, `${1}_${2}`)
	snake = strings.ToLower(snake)
	if !snakeValidRegex.MatchString(snake) {
		return ""
	}
	return snake
}
