package relationDB

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
