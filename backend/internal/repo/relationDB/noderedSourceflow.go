package relationDB

import (
	"backend/internal/types"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GroupedSourceFlowItem 分组和未分组source flow统一返回结构
type GroupedSourceFlowItem struct {
	ID          string    `gorm:"column:id" json:"id"`                    // ID
	Category    string    `gorm:"column:category" json:"category"`        // 分类 group-分组 file-文件
	GroupType   int64     `gorm:"column:group_type" json:"groupType"`     // 类型：GROUP 表示分组，BIZ 表示未分组的source flow
	Name        string    `gorm:"column:name" json:"name"`                // 名称
	Description string    `gorm:"column:description" json:"description"`  // 描述
	GroupID     *int64    `gorm:"column:group_id" json:"groupId"`         // 分组ID
	Sort        int32     `gorm:"column:sort" json:"sort"`                // 排序字段
	CreateAt    time.Time `gorm:"column:create_at" json:"createAt"`       // 创建时间
	Creator     string    `gorm:"column:creator" json:"creator"`          // 创建人
	FlowName    string    `gorm:"column:flow_name" json:"flowName"`       // flow名称
	FlowID      string    `gorm:"column:flow_id" json:"flowId"`           // flow ID
	FlowStatus  string    `gorm:"column:flow_status" json:"flowStatus"`   // flow状态
	Template    string    `gorm:"column:template" json:"template"`        // 模板类型
	HasChildren bool      `gorm:"column:has_children" json:"hasChildren"` // 是否有子节点
}

// GetGroupedSourceFlowList 按分组获取source flow列表
// GetGroupedFlowList 按分组获取flow列表（内部通用方法）
func (m *NoderedSourceFlowRepo) GetGroupedFlowList(
	db *gorm.DB,
	req *types.GroupPageRequest,
	template string,
) ([]*GroupedSourceFlowItem, int64, error) {

	var (
		items []*GroupedSourceFlowItem
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
		    null as flow_name,
		    null as flow_id,
		    null as flow_status,
		    null as template,
		    EXISTS (
				SELECT 1
				FROM supos_node_flows u
				WHERE u.group_id = g.id
    		) AS has_children
		FROM resource_group g
		WHERE g.type = $1
		UNION ALL
		SELECT
		    u.id::varchar AS id,
		 	'file' AS category,
		    $1 AS group_type,
		    u.flow_name AS name,
		    u.description AS description,
		    u.group_id AS group_id,
		    COALESCE(r.mark,0) AS sort,
		    u.create_time AS create_at,
		    u.creator AS creator,
		    u.flow_name AS flow_name,
		    u.flow_id AS flow_id,
		    u.flow_status AS flow_status,
		    u.template AS template,
		    FALSE AS has_children
		FROM supos_node_flows u
		LEFT JOIN supos_node_flow_top_recodes r
		    ON u.id = r.id
		WHERE u.template = $2
	`

	// $1 -> groupType, $2 -> template
	args = append(args, req.GroupType, template)
	paramIdx := 3

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
		logx.Errorf("failed to count grouped flow list: %v", err)
		return nil, 0, err
	}

	// ---------- pagination ----------
	offset := (req.PageNo - 1) * req.PageSize
	querySQL += fmt.Sprintf(
		" ORDER BY sort DESC, create_at DESC LIMIT $%d OFFSET $%d",
		paramIdx, paramIdx+1,
	)
	args = append(args, req.PageSize, offset)

	logx.Debugf("group flow sql: %s , args: %+v", querySQL, args)

	// ---------- query ----------
	if err := db.Raw(querySQL, args...).Scan(&items).Error; err != nil {
		logx.Errorf("failed to get grouped flow list: %v", err)
		return nil, 0, err
	}

	return items, total, nil
}

// GetGroupedSourceFlowList 按分组获取source flow列表
func (m *NoderedSourceFlowRepo) GetGroupedSourceFlowList(db *gorm.DB, req *types.GroupPageRequest) ([]*GroupedSourceFlowItem, int64, error) {
	// 根据 groupType 判断 template
	var template string
	switch req.GroupType {
	case 1:
		template = "node-red"
	case 2:
		template = "event-flow"
	default:
		return nil, 0, fmt.Errorf("invalid group type: %d", req.GroupType)
	}

	return m.GetGroupedFlowList(db, req, template)
}

// GetGroupedEventFlowList 按分组获取event flow列表
func (m *NoderedSourceFlowRepo) GetGroupedEventFlowList(db *gorm.DB, req *types.GroupPageRequest) ([]*GroupedSourceFlowItem, int64, error) {
	// 根据 groupType 判断 template
	var template string
	switch req.GroupType {
	case 1:
		template = "node-red"
	case 2:
		template = "event-flow"
	default:
		return nil, 0, fmt.Errorf("invalid group type: %d", req.GroupType)
	}

	return m.GetGroupedFlowList(db, req, template)
}

/*
这个是参考样例
使用教程:
1. 将NoderedSourceFlow全局替换为模型的表名
2. 完善todo
*/

type NoderedSourceFlowRepo struct {
	db *gorm.DB
}

func NewNoderedSourceFlowRepo(in context.Context) *NoderedSourceFlowRepo {
	return &NoderedSourceFlowRepo{db: GetDb(in)}
}

type NoderedSourceFlowFilter struct {
	//todo 添加过滤字段
	ID        int64
	Name      string
	NameLike  string
	Template  string
	Templates []string
	// FlowType int32
	FlowID string
}

func (p NoderedSourceFlowRepo) fmtFilter(ctx context.Context, f NoderedSourceFlowFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if len(f.Templates) > 0 {
		db = db.Where("template IN ?", f.Templates)
	}
	if f.Template != "" {
		db = db.Where("template = ?", f.Template)
	}
	if f.Name != "" {
		db = db.Where("flow_name = ?", f.Name)
	}
	if f.NameLike != "" {
		db = db.Where("flow_name LIKE ?", "%"+f.NameLike+"%")
	}
	if f.FlowID != "" {
		db = db.Where("flow_id = ?", f.FlowID)
	}
	return db
}

func (p NoderedSourceFlowRepo) Insert(ctx context.Context, data *NoderedSourceFlow) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NoderedSourceFlowRepo) FindOneByFilter(ctx context.Context, f NoderedSourceFlowFilter) (*NoderedSourceFlow, error) {
	var result NoderedSourceFlow
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NoderedSourceFlowRepo) FindByFilter(ctx context.Context, f NoderedSourceFlowFilter, page *stores.PageInfo) ([]*NoderedSourceFlow, error) {
	var results []*NoderedSourceFlow
	db := p.fmtFilter(ctx, f).Model(&NoderedSourceFlow{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NoderedSourceFlowRepo) CountByFilter(ctx context.Context, f NoderedSourceFlowFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&NoderedSourceFlow{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NoderedSourceFlowRepo) Update(ctx context.Context, data *NoderedSourceFlow) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NoderedSourceFlowRepo) DeleteByFilter(ctx context.Context, f NoderedSourceFlowFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&NoderedSourceFlow{}).Error
	return stores.ErrFmt(err)
}

func (p NoderedSourceFlowRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&NoderedSourceFlow{}).Error
	return stores.ErrFmt(err)
}
func (p NoderedSourceFlowRepo) FindOne(ctx context.Context, id int64) (*NoderedSourceFlow, error) {
	var result NoderedSourceFlow
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NoderedSourceFlowRepo) MultiInsert(ctx context.Context, data []*NoderedSourceFlow) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&NoderedSourceFlow{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d NoderedSourceFlowRepo) UpdateWithField(ctx context.Context, f NoderedSourceFlowFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&NoderedSourceFlow{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

// 关系替换：按 flow 覆盖写 model 关联
func (r NoderedSourceFlowRepo) ReplaceModels(ctx context.Context, parentID int64, modelAlias []string) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("parent_id = ?", parentID).Delete(&NoderedSourceFlowNode{}).Error; err != nil {
		tx.Rollback()
		return stores.ErrFmt(err)
	}
	if len(modelAlias) > 0 {
		recs := make([]*NoderedSourceFlowNode, 0, len(modelAlias))
		for _, alias := range modelAlias {
			recs = append(recs, &NoderedSourceFlowNode{ParentID: parentID, Alias: alias})
		}
		if err := tx.Create(&recs).Error; err != nil {
			tx.Rollback()
			return stores.ErrFmt(err)
		}
	}
	return stores.ErrFmt(tx.Commit().Error)
}

// SelectByModelIDs 根据模型ID集合查询关联的 Flow 列表
func (r NoderedSourceFlowRepo) SelectByModelIDs(ctx context.Context, modelIDs []int64) ([]*NoderedSourceFlow, error) {
	if len(modelIDs) == 0 {
		return []*NoderedSourceFlow{}, nil
	}
	var parentIds []int64
	if err := r.db.WithContext(ctx).Model(&NoderedSourceFlowNode{}).Where("node_id IN ?", modelIDs).Pluck("parent_id", &parentIds).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	if len(parentIds) == 0 {
		return []*NoderedSourceFlow{}, nil
	}
	var flows []*NoderedSourceFlow
	if err := r.db.WithContext(ctx).Where("id IN ?", parentIds).Find(&flows).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	return flows, nil
}

// FindLatestByNodeID returns the latest associated flow for a given UNS node id.
func (r NoderedSourceFlowRepo) FindLatestByNodeID(ctx context.Context, nodeID int64) (*NoderedSourceFlow, error) {
	// Step 1: query latest relation via model (respects naming strategy)
	var rel NoderedSourceFlowNode
	q := r.db.WithContext(ctx).Model(&NoderedSourceFlowNode{}).
		Where("node_id = ?", nodeID).
		Order("created_time DESC").
		Limit(1).
		Take(&rel)
	if err := q.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, stores.ErrFmt(err)
	}
	// Step 2: fetch flow by parent id
	var flow NoderedSourceFlow
	if err := r.db.WithContext(ctx).Where("id = ?", rel.ParentID).First(&flow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, stores.ErrFmt(err)
	}
	return &flow, nil
}

// FindLatestByAlias returns the latest associated flow for a given UNS alias.
func (r NoderedSourceFlowRepo) FindLatestByAlias(ctx context.Context, alias string) (*NoderedSourceFlow, error) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return nil, nil
	}
	var rel NoderedSourceFlowNode
	q := r.db.WithContext(ctx).Model(&NoderedSourceFlowNode{}).
		Where("alias = ?", alias).
		Order("create_time DESC").
		Limit(1).
		Take(&rel)
	if err := q.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, stores.ErrFmt(err)
	}
	var flow NoderedSourceFlow
	if err := r.db.WithContext(ctx).Where("id = ?", rel.ParentID).First(&flow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, stores.ErrFmt(err)
	}
	return &flow, nil
}

// FindAvailableFlowName ensures flow_name uniqueness by appending -N suffix when needed, scoped by template(flow type).
func (r NoderedSourceFlowRepo) FindAvailableFlowName(ctx context.Context, base string, flowType string) (string, int, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		return "", 0, fmt.Errorf("flow name empty")
	}
	var rows []NoderedSourceFlow
	like := base + "%"
	db := r.db.WithContext(ctx).Where("flow_name LIKE ?", like)
	if strings.TrimSpace(flowType) != "" {
		db = db.Where("template = ?", strings.TrimSpace(flowType))
	}
	if err := db.Find(&rows).Error; err != nil {
		return "", 0, stores.ErrFmt(err)
	}
	suffixRe := regexp.MustCompile(`^(.*?)-(\d+)$`)
	maxN := 0
	existsBase := false
	for _, r0 := range rows {
		if r0.FlowName == base {
			existsBase = true
			continue
		}
		if m := suffixRe.FindStringSubmatch(r0.FlowName); len(m) == 3 && m[1] == base {
			var n int
			fmt.Sscanf(m[2], "%d", &n)
			if n > maxN {
				maxN = n
			}
		}
	}
	if !existsBase {
		return base, 0, nil
	}
	return fmt.Sprintf("%s(%d)", base, maxN+1), maxN + 1, nil
}

// SelectByAliases returns flows associated with any of the given aliases.
func (r NoderedSourceFlowRepo) SelectByAliases(ctx context.Context, aliases []string) ([]*NoderedSourceFlow, error) {
	clean := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		if v := strings.TrimSpace(alias); v != "" {
			clean = append(clean, v)
		}
	}
	if len(clean) == 0 {
		return []*NoderedSourceFlow{}, nil
	}
	var parentIDs []int64
	if err := r.db.WithContext(ctx).
		Model(&NoderedSourceFlowNode{}).
		Where("alias IN ?", clean).
		Pluck("parent_id", &parentIDs).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	if len(parentIDs) == 0 {
		return []*NoderedSourceFlow{}, nil
	}
	var flows []*NoderedSourceFlow
	if err := r.db.WithContext(ctx).
		Where("id IN ?", parentIDs).
		Find(&flows).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	return flows, nil
}
