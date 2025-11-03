package dao

import (
	"backend/internal/logic/supos/uns/dashboard/model"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// DashboardMarkedMapper Dashboard 置顶标记数据访问对象
type DashboardMarkedMapper struct {
	conn   sqlx.SqlConn
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardMarkedMapper 创建 DashboardMarkedMapper 实例
func NewDashboardMarkedMapper(conn sqlx.SqlConn, ctx context.Context) *DashboardMarkedMapper {
	return &DashboardMarkedMapper{
		conn:   conn,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// Insert 插入置顶标记
func (m *DashboardMarkedMapper) Insert(mark *model.DashboardMarkModel) error {
	query := `INSERT INTO uns_dashboard_top_recodes (id, user_id) VALUES ($1, $2)`

	_, err := m.conn.ExecCtx(m.ctx, query, mark.ID, mark.UserID)
	if err != nil {
		m.logger.Errorf("failed to insert dashboard mark: %v", err)
		return err
	}
	return nil
}

// Delete 删除置顶标记
func (m *DashboardMarkedMapper) Delete(id string, userID string) error {
	query := `DELETE FROM uns_dashboard_top_recodes WHERE id = $1 AND user_id = $2`

	_, err := m.conn.ExecCtx(m.ctx, query, id, userID)
	if err != nil {
		m.logger.Errorf("failed to delete dashboard mark: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 Dashboard ID 删除所有置顶标记
func (m *DashboardMarkedMapper) DeleteById(id string) error {
	query := `DELETE FROM uns_dashboard_top_recodes WHERE id = $1`

	_, err := m.conn.ExecCtx(m.ctx, query, id)
	if err != nil {
		m.logger.Errorf("failed to delete dashboard mark by id: %v", err)
		return err
	}
	return nil
}
