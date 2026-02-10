package relationDB

import (
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DashboardMapper Dashboard 数据访问对象
type DashboardMapper struct {
}

// SelectById 根据 ID 查询 Dashboard
func (m *DashboardMapper) SelectById(db *gorm.DB, id string) (*DashboardModel, error) {
	var dashboard DashboardModel
	err := db.Where("id = ?", id).First(&dashboard).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logx.Errorf("failed to select dashboard by id: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// Insert 插入 Dashboard
func (m *DashboardMapper) Insert(db *gorm.DB, dashboard *DashboardModel) error {
	err := db.Create(dashboard).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard: %v", err)
		return err
	}
	return nil
}
func (m *DashboardMapper) SaveBatch(db *gorm.DB, dashboard []*DashboardModel) error {
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(dashboard).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard: %v", err)
		return err
	}
	return nil
}
func (m *DashboardMapper) Save(ctx context.Context, dashboard []*DashboardModel) error {
	db := GetDb(ctx)
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Save(dashboard).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard: %v", err)
		return err
	}
	return nil
}

// UpdateById 根据 ID 更新 Dashboard
func (m *DashboardMapper) UpdateById(db *gorm.DB, dashboard *DashboardModel) error {
	// 使用 map 更新非零值字段，避免gorm默认的“忽略零值”行为
	// 这里假设所有字段都需要更新
	err := db.Model(&DashboardModel{}).Where("id = ?", dashboard.ID).Updates(dashboard).Error
	if err != nil {
		logx.Errorf("failed to update dashboard: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 ID 删除 Dashboard
func (m *DashboardMapper) DeleteById(db *gorm.DB, id string) error {
	err := db.Where("id = ?", id).Delete(&DashboardModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard: %v", err)
		return err
	}
	return nil
}

// SelectByFlowNames 根据名称列表查询 Dashboard
func (m *DashboardMapper) SelectByFlowNames(db *gorm.DB, names []string) ([]*DashboardModel, error) {
	if len(names) == 0 {
		return []*DashboardModel{}, nil
	}
	var dashboards []*DashboardModel
	err := db.Where("name IN ?", names).Find(&dashboards).Error
	if err != nil {
		logx.Errorf("failed to select dashboards by names: %v", err)
		return nil, err
	}
	return dashboards, nil
}
func (m *DashboardMapper) SelectByNameAndType(db *gorm.DB, name string, dashType int) ([]*DashboardModel, error) {
	var dashboards []*DashboardModel
	err := db.Where("name = ?", name).Where("type=?", dashType).Find(&dashboards).Error
	if err != nil {
		logx.Errorf("failed to select dashboards by NameAndType: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SaveOrIgnoreBatch 批量保存或忽略
func (m *DashboardMapper) SaveOrIgnoreBatch(db *gorm.DB, dashboards []*DashboardModel) error {
	if len(dashboards) == 0 {
		return nil
	}
	// GORM v2 的 Clauses(clause.OnConflict{DoNothing: true}) 提供了优雅的方式
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&dashboards).Error
	if err != nil {
		logx.Errorf("failed to batch save or ignore dashboards: %v", err)
		return err
	}
	return nil
}

// DashboardExtends Dashboard 扩展信息（包含置顶标记）
type DashboardExtends struct {
	DashboardModel
	Mark     *int       `db:"mark" json:"mark,omitzero"`          // 置顶标记
	MarkTime *time.Time `db:"mark_time" json:"markTime,omitzero"` // 置顶时间
}

// SelectDashboard 分页查询 Dashboard（包含置顶信息）
func (m *DashboardMapper) SelectDashboard(
	db *gorm.DB,
	userID string,
	fuzzyName string,
	typ *int,
	orderCode string,
	asc bool,
	pageNo int64,
	pageSize int64,
	countTotal *int64,
) ([]*DashboardExtends, error) {
	var dashboards []*DashboardExtends
	query := db.
		Table("uns_dashboard a").
		Select("a.*, b.mark, b.mark_time").
		Joins("LEFT JOIN uns_dashboard_top_recodes b ON a.id = b.id AND b.user_id = ?", userID)

	if fuzzyName != "" {
		searchPattern := "%" + escapeLikePattern(fuzzyName) + "%"
		query = query.Where("(a.name LIKE ? OR a.description LIKE ?)", searchPattern, searchPattern)
	}

	if typ != nil {
		query = query.Where("a.type = ?", *typ)
	}

	if countTotal != nil {
		er := query.Count(countTotal).Error
		if er != nil {
			return nil, er
		}
	}

	// 排序
	orders := base.StringBuilder{}
	orders.Grow(64)
	orders.Append("b.mark asc, ")
	if orderCode == "" {
		orders.Append(" b.mark_time desc, a.create_time desc ")
	} else {
		orders.Append(fmt.Sprintf("%s %s", orderCode, base.SanYuan(asc, "ASC", "DESC")))
	}
	query = query.Order(orders.String())
	// 分页
	offset := (pageNo - 1) * pageSize
	query = query.Limit(int(pageSize)).Offset(int(offset))
	err := query.Find(&dashboards).Error
	if err != nil {
		logx.Errorf("failed to select dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectDashboardCount 查询 Dashboard 总数
func (m *DashboardMapper) SelectDashboardCount(db *gorm.DB, fuzzyName string, typ *int) (int64, error) {
	var count int64
	query := db.Model(&DashboardModel{})

	if fuzzyName != "" {
		searchPattern := "%" + fuzzyName + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	if typ != nil {
		query = query.Where("type = ?", *typ)
	}

	err := query.Count(&count).Error
	if err != nil {
		logx.Errorf("failed to count dashboards: %v", err)
		return 0, err
	}
	return count, nil
}

// SelectAll selects all DashboardModel from the database.
func (m *DashboardMapper) SelectAll(db *gorm.DB) ([]*DashboardModel, error) {
	var dashboards []*DashboardModel
	err := db.Find(&dashboards).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*DashboardModel{}, nil
		}
		logx.Errorf("failed to select all dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectByIds selects multiple DashboardModel from the database by their IDs.
func (m *DashboardMapper) SelectByIds(db *gorm.DB, ids []string) ([]*DashboardModel, error) {
	if len(ids) == 0 {
		return []*DashboardModel{}, nil
	}
	var dashboards []*DashboardModel
	err := db.Where("id IN ?", ids).Find(&dashboards).Error
	if err != nil {
		logx.Errorf("failed to select dashboards by ids: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// DeleteBatchIds deletes multiple dashboards from the database by their IDs.
func (m *DashboardMapper) DeleteBatchIds(db *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	err := db.Where("id IN ?", ids).Delete(&DashboardModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboards by ids: %v", err)
		return err
	}
	return nil
}

// GroupedDashboardItem 分组和未分组dashboard统一返回结构
type GroupedDashboardItem struct {
	ID          string    `gorm:"column:id" json:"id"`                    // ID
	Category    string    `gorm:"column:category" json:"category"`        // 分类 group-分组 file-文件
	GroupType   int64     `gorm:"column:group_type" json:"groupType"`     // 类型：1-sourceflow 2-eventflow 3-datasource
	Name        string    `gorm:"column:name" json:"name"`                // 名称
	Description string    `gorm:"column:description" json:"description"`  // 描述
	GroupID     *int64    `gorm:"column:group_id" json:"groupId"`         // 分组ID
	Sort        int32     `gorm:"column:sort" json:"sort"`                // 排序字段
	CreateAt    time.Time `gorm:"column:create_at" json:"createAt"`       // 创建时间
	Creator     string    `gorm:"column:creator" json:"creator"`          // 创建人
	HasChildren bool      `gorm:"column:has_children" json:"hasChildren"` // 是否有子节点
}

// GetGroupedDashboardList 按分组获取dashboard列表
func (m *DashboardMapper) GetGroupedDashboardList(
	db *gorm.DB,
	req *types.GroupPageRequest,
) ([]*GroupedDashboardItem, int64, error) {

	var (
		items []*GroupedDashboardItem
		total int64
	)

	whereConditions := make([]string, 0)
	args := make([]interface{}, 0)

	// ---------- base SQL ----------
	baseSQL := `
		SELECT
		    g.id::varchar AS id,
		    'group' AS category,
		    g.type AS group_type,
		    g.name AS name,
		    g.description AS description,
		    NULL AS group_id,
		    g.sort AS sort,
		    g.create_at AS create_at,
		    g.creator AS creator,
		    EXISTS (
				SELECT 1
				FROM uns_dashboard u
				WHERE u.group_id = g.id
   			) AS has_children
		FROM resource_group g where g.type = 3
		UNION ALL
		SELECT
		    u.id AS id,
		    'file' AS category,
		    $1 AS group_type,
		    u.name AS name,
		    u.description AS description,
		    u.group_id AS group_id,
		    COALESCE(r.mark,0) AS sort,
		    u.create_time AS create_at,
			u.creator AS creator,
			false AS has_children
		FROM uns_dashboard u
		LEFT JOIN uns_dashboard_top_recodes r
		    ON u.id = r.id::varchar
	`

	// $1
	args = append(args, req.GroupType)
	paramIdx := 2

	// ---------- WHERE conditions ----------

	// groupId
	if req.GroupId > 0 {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf("group_id = $%d", paramIdx),
		)
		args = append(args, req.GroupId)
		paramIdx++
	} else {
		whereConditions = append(whereConditions, "group_id IS NULL")
	}

	// category
	if req.Category != "" {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf("category = $%d", paramIdx),
		)
		args = append(args, req.Category)
		paramIdx++
	}

	// keyword
	if req.K != "" {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"(name ILIKE $%d OR description ILIKE $%d)",
				paramIdx, paramIdx+1,
			),
		)
		like := "%" + req.K + "%"
		args = append(args, like, like)
		paramIdx += 2
	}

	// ---------- final query SQL ----------
	querySQL := baseSQL
	if len(whereConditions) > 0 {
		querySQL = `
			SELECT * FROM (` + baseSQL + `) t
			WHERE ` + strings.Join(whereConditions, " AND ")
	}

	// ---------- COUNT ----------
	countSQL := "SELECT COUNT(*) FROM (" + querySQL + ") t"
	if err := db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		logx.Errorf("failed to count grouped dashboard list: %v", err)
		return nil, 0, err
	}

	// ---------- sorting ----------
	sortSQL := " ORDER BY sort DESC"
	if req.OrderCode != "" {
		// 如果有指定排序字段，替换默认的 create_at 字段
		if req.OrderCode == "createAt" {
			req.OrderCode = "create_at"
		}

		ascDesc := "DESC"
		if req.IsAsc {
			ascDesc = "ASC"
		}
		sortSQL += fmt.Sprintf(", %s %s", req.OrderCode, ascDesc)
	} else {
		// 默认使用 create_at 降序
		sortSQL += ", create_at DESC"
	}

	// ---------- pagination ----------
	offset := (req.PageNo - 1) * req.PageSize
	querySQL += fmt.Sprintf(
		"%s LIMIT $%d OFFSET $%d",
		sortSQL, paramIdx, paramIdx+1,
	)
	args = append(args, req.PageSize, offset)

	logx.Debugf("group sql: %s , args: %+v", querySQL, args)

	// ---------- query ----------
	if err := db.Raw(querySQL, args...).Scan(&items).Error; err != nil {
		logx.Errorf("failed to get grouped dashboard list: %v", err)
		return nil, 0, err
	}

	return items, total, nil
}

// SelectDashboardsToInit selects dashboards that need to be initialized.
func (m *DashboardMapper) SelectDashboardsToInit(db *gorm.DB) ([]*DashboardModel, error) {
	var dashboards []*DashboardModel
	err := db.
		Where("need_init = ? AND type = ? AND json_content IS NOT NULL AND json_content != ?", true, 1, "").
		Find(&dashboards).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*DashboardModel{}, nil
		}
		logx.Errorf("failed to select dashboards to init: %v", err)
		return nil, err
	}
	return dashboards, nil
}
