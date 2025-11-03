package dao

import (
	"backend/internal/logic/supos/uns/dashboard/model"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DashboardMapper Dashboard 数据访问对象
type DashboardMapper struct {
	conn   sqlx.SqlConn
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardMapper 创建 DashboardMapper 实例
func NewDashboardMapper(conn sqlx.SqlConn, ctx context.Context) *DashboardMapper {
	return &DashboardMapper{
		conn:   conn,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// SelectById 根据 ID 查询 Dashboard
func (m *DashboardMapper) SelectById(id string) (*model.DashboardModel, error) {
	query := `SELECT id, name, type, need_init, description, json_content, creator, update_time, create_time 
              FROM uns_dashboard WHERE id = $1`

	var dashboard model.DashboardModel
	err := m.conn.QueryRowCtx(m.ctx, &dashboard, query, id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to select dashboard by id: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// Insert 插入 Dashboard
func (m *DashboardMapper) Insert(dashboard *model.DashboardModel) error {
	query := `INSERT INTO uns_dashboard (id, name, type, need_init, description, json_content, creator, update_time, create_time) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := m.conn.ExecCtx(m.ctx, query,
		dashboard.ID,
		dashboard.Name,
		dashboard.Type,
		dashboard.NeedInit,
		dashboard.Description,
		dashboard.JsonContent,
		dashboard.Creator,
	)
	if err != nil {
		m.logger.Errorf("failed to insert dashboard: %v", err)
		return err
	}
	return nil
}

// UpdateById 根据 ID 更新 Dashboard
func (m *DashboardMapper) UpdateById(dashboard *model.DashboardModel) error {
	query := `UPDATE uns_dashboard SET name = $1, type = $2, need_init = $3, description = $4, 
              json_content = $5, update_time = CURRENT_TIMESTAMP WHERE id = $6`

	_, err := m.conn.ExecCtx(m.ctx, query,
		dashboard.Name,
		dashboard.Type,
		dashboard.NeedInit,
		dashboard.Description,
		dashboard.JsonContent,
		dashboard.ID,
	)
	if err != nil {
		m.logger.Errorf("failed to update dashboard: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 ID 删除 Dashboard
func (m *DashboardMapper) DeleteById(id string) error {
	query := `DELETE FROM uns_dashboard WHERE id = $1`

	_, err := m.conn.ExecCtx(m.ctx, query, id)
	if err != nil {
		m.logger.Errorf("failed to delete dashboard: %v", err)
		return err
	}
	return nil
}

// SelectByFlowNames 根据名称列表查询 Dashboard
func (m *DashboardMapper) SelectByFlowNames(names []string) ([]*model.DashboardModel, error) {
	if len(names) == 0 {
		return []*model.DashboardModel{}, nil
	}

	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	for i, name := range names {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(`SELECT id, name, type, need_init, description, json_content, creator, update_time, create_time 
                          FROM uns_dashboard WHERE name IN (%s)`, strings.Join(placeholders, ","))

	var dashboards []*model.DashboardModel
	err := m.conn.QueryRowsCtx(m.ctx, &dashboards, query, args...)
	if err != nil {
		m.logger.Errorf("failed to select dashboards by names: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SaveOrIgnoreBatch 批量保存或忽略
func (m *DashboardMapper) SaveOrIgnoreBatch(dashboards []*model.DashboardModel) error {
	if len(dashboards) == 0 {
		return nil
	}

	// 构建批量插入 SQL
	valuePlaceholders := make([]string, len(dashboards))
	args := make([]interface{}, 0, len(dashboards)*5)
	argIndex := 1

	for i, db := range dashboards {
		valuePlaceholders[i] = fmt.Sprintf("($%d, $%d, $%d, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			argIndex, argIndex+1, argIndex+2)
		args = append(args, db.ID, db.Name, db.Description)
		argIndex += 3
	}

	query := fmt.Sprintf(`INSERT INTO uns_dashboard (id, name, description, update_time, create_time) 
                          VALUES %s ON CONFLICT (id) DO NOTHING`, strings.Join(valuePlaceholders, ","))

	_, err := m.conn.ExecCtx(m.ctx, query, args...)
	if err != nil {
		m.logger.Errorf("failed to batch save or ignore dashboards: %v", err)
		return err
	}
	return nil
}

// DashboardExtends Dashboard 扩展信息（包含置顶标记）
type DashboardExtends struct {
	model.DashboardModel
	Mark     *int   `db:"mark" json:"mark,omitzero"`          // 置顶标记
	MarkTime *int64 `db:"mark_time" json:"markTime,omitzero"` // 置顶时间
}

// SelectDashboard 分页查询 Dashboard（包含置顶信息）
func (m *DashboardMapper) SelectDashboard(
	userID string,
	fuzzyName string,
	typ *int,
	orderCode string,
	descOrAsc string,
	pageNo int64,
	pageSize int64,
) ([]*DashboardExtends, error) {
	// 构建动态 SQL
	var conditions []string
	args := []interface{}{userID}
	argIndex := 2

	baseQuery := `SELECT a.*, b.mark, b.mark_time 
                  FROM uns_dashboard a 
                  LEFT JOIN uns_dashboard_top_recodes b ON a.id = b.id AND b.user_id = $1 
                  WHERE 1=1`

	if fuzzyName != "" {
		conditions = append(conditions, fmt.Sprintf(" AND (a.name LIKE $%d OR a.description LIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+fuzzyName+"%")
		argIndex++
	}

	if typ != nil {
		conditions = append(conditions, fmt.Sprintf(" AND a.type = $%d", argIndex))
		args = append(args, *typ)
		argIndex++
	}

	query := baseQuery + strings.Join(conditions, "")

	// 排序
	if orderCode == "" {
		query += " ORDER BY b.mark ASC, b.mark_time DESC, a.create_time DESC"
	} else {
		query += fmt.Sprintf(" ORDER BY %s %s", orderCode, descOrAsc)
	}

	// 分页
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (pageNo-1)*pageSize)

	var dashboards []*DashboardExtends
	err := m.conn.QueryRowsCtx(m.ctx, &dashboards, query, args...)
	if err != nil {
		m.logger.Errorf("failed to select dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectDashboardCount 查询 Dashboard 总数
func (m *DashboardMapper) SelectDashboardCount(fuzzyName string, typ *int) (int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	query := `SELECT COUNT(*) FROM uns_dashboard WHERE 1=1`

	if fuzzyName != "" {
		conditions = append(conditions, fmt.Sprintf(" AND (name LIKE $%d OR description LIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+fuzzyName+"%")
		argIndex++
	}

	if typ != nil {
		conditions = append(conditions, fmt.Sprintf(" AND type = $%d", argIndex))
		args = append(args, *typ)
		argIndex++
	}

	query += strings.Join(conditions, "")

	var count int64
	err := m.conn.QueryRowCtx(m.ctx, &count, query, args...)
	if err != nil {
		m.logger.Errorf("failed to count dashboards: %v", err)
		return 0, err
	}
	return count, nil
}

// SelectAll selects all DashboardModel from the database.
func (m *DashboardMapper) SelectAll() ([]*model.DashboardModel, error) {
	var dashboards []*model.DashboardModel
	query := "SELECT id, name, type, need_init, description, json_content, creator, update_time, create_time FROM uns_dashboard"
	err := m.conn.QueryRowsCtx(m.ctx, &dashboards, query)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return []*model.DashboardModel{}, nil
		}
		m.logger.Errorf("failed to select all dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectByIds selects multiple DashboardModel from the database by their IDs.
func (m *DashboardMapper) SelectByIds(ids []string) ([]*model.DashboardModel, error) {
	if len(ids) == 0 {
		return []*model.DashboardModel{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, name, type, need_init, description, json_content, creator, update_time, create_time 
                          FROM uns_dashboard WHERE id IN (%s)`, strings.Join(placeholders, ","))

	var dashboards []*model.DashboardModel
	err := m.conn.QueryRowsCtx(m.ctx, &dashboards, query, args...)
	if err != nil {
		m.logger.Errorf("failed to select dashboards by ids: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// DeleteBatchIds deletes multiple dashboards from the database by their IDs.
func (m *DashboardMapper) DeleteBatchIds(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`DELETE FROM uns_dashboard WHERE id IN (%s)`, strings.Join(placeholders, ","))

	_, err := m.conn.ExecCtx(m.ctx, query, args...)
	if err != nil {
		m.logger.Errorf("failed to delete dashboards by ids: %v", err)
		return err
	}
	return nil
}

// SelectDashboardsToInit selects dashboards that need to be initialized.
func (m *DashboardMapper) SelectDashboardsToInit() ([]*model.DashboardModel, error) {
	query := `SELECT id, name, type, need_init, description, json_content, creator, update_time, create_time 
              FROM uns_dashboard WHERE need_init = true AND type = 1 AND json_content IS NOT NULL AND json_content != ''`

	var dashboards []*model.DashboardModel
	err := m.conn.QueryRowsCtx(m.ctx, &dashboards, query)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return []*model.DashboardModel{}, nil
		}
		m.logger.Errorf("failed to select dashboards to init: %v", err)
		return nil, err
	}
	return dashboards, nil
}
