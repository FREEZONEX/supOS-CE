package dao

import (
	"backend/internal/logic/supos/uns/dashboard/model"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DashboardRefMapper Dashboard 引用关系数据访问对象
type DashboardRefMapper struct {
	conn   sqlx.SqlConn
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardRefMapper 创建 DashboardRefMapper 实例
func NewDashboardRefMapper(conn sqlx.SqlConn, ctx context.Context) *DashboardRefMapper {
	return &DashboardRefMapper{
		conn:   conn,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// Insert 插入 Dashboard 引用关系
func (m *DashboardRefMapper) Insert(ref *model.DashboardRefModel) error {
	query := `INSERT INTO uns_dashboard_ref (dashboard_id, uns_alias, create_at) 
              VALUES ($1, $2, CURRENT_TIMESTAMP)`

	_, err := m.conn.ExecCtx(m.ctx, query, ref.DashboardID, ref.UnsAlias)
	if err != nil {
		m.logger.Errorf("failed to insert dashboard ref: %v", err)
		return err
	}
	return nil
}

// DeleteByDashboardId 根据 Dashboard ID 删除引用关系
func (m *DashboardRefMapper) DeleteByDashboardId(dashboardID string) error {
	query := `DELETE FROM uns_dashboard_ref WHERE dashboard_id = $1`

	_, err := m.conn.ExecCtx(m.ctx, query, dashboardID)
	if err != nil {
		m.logger.Errorf("failed to delete dashboard ref: %v", err)
		return err
	}
	return nil
}

// GetByUns 根据 UNS 别名获取 Dashboard
func (m *DashboardRefMapper) GetByUns(unsAlias string) (*model.DashboardModel, error) {
	query := `SELECT a.id, a.name, a.type, a.need_init, a.description, a.json_content, a.creator, a.update_time, a.create_time 
              FROM uns_dashboard a 
              LEFT JOIN uns_dashboard_ref b ON a.id = b.dashboard_id 
              WHERE b.uns_alias = $1 LIMIT 1`

	var dashboard model.DashboardModel
	err := m.conn.QueryRowCtx(m.ctx, &dashboard, query, unsAlias)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to get dashboard by uns: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// SelectByUnsAlias 根据 UNS 别名查询引用关系
func (m *DashboardRefMapper) SelectByUnsAlias(unsAlias string) (*model.DashboardRefModel, error) {
	query := `SELECT dashboard_id, uns_alias, create_at FROM uns_dashboard_ref WHERE uns_alias = $1`

	var ref model.DashboardRefModel
	err := m.conn.QueryRowCtx(m.ctx, &ref, query, unsAlias)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to select dashboard ref: %v", err)
		return nil, err
	}
	return &ref, nil
}

// SelectByUnsAliases selects dashboard references by a list of UNS aliases.
func (m *DashboardRefMapper) SelectByUnsAliases(aliases []string) ([]*model.DashboardRefModel, error) {
	if len(aliases) == 0 {
		return []*model.DashboardRefModel{}, nil
	}

	placeholders := make([]string, len(aliases))
	args := make([]interface{}, len(aliases))
	for i, alias := range aliases {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = alias
	}

	query := fmt.Sprintf(`SELECT dashboard_id, uns_alias, create_at FROM uns_dashboard_ref WHERE uns_alias IN (%s)`, strings.Join(placeholders, ","))

	var refs []*model.DashboardRefModel
	err := m.conn.QueryRowsCtx(m.ctx, &refs, query, args...)
	if err != nil {
		m.logger.Errorf("failed to select dashboard refs by aliases: %v", err)
		return nil, err
	}
	return refs, nil
}
