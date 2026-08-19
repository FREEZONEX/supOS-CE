package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidUnsNode     = errors.New("invalid uns node")
	ErrUnsRestoreConflict = errors.New("uns restore namespace conflict")
	ErrUnsNotInRecycle    = errors.New("uns node is not in recycle")
)

type unsLabelNode struct {
	LabelID int64 `gorm:"column:label_id;primaryKey"`
	NodeID  int64 `gorm:"column:node_id;primaryKey"`
	CreateTime
}

func (unsLabelNode) TableName() string { return "uns_namespace_label_nodeid" }

func loadUnsNode(db *gorm.DB, id int64, includeDeleted bool) (UnsNode, error) {
	q := db
	if includeDeleted {
		q = q.Unscoped()
	} else {
		q = q.Where("deleted_time = 0")
	}
	q = q.Where("id = ?", id)
	var node UnsNode
	err := q.Take(&node).Error
	return node, err
}

func listUnsSubtree(db *gorm.DB, namespace string, includeDeleted bool) ([]UnsNode, error) {
	q := db
	if includeDeleted {
		q = q.Unscoped()
	} else {
		q = q.Where("deleted_time = 0")
	}
	q = q.Where("(namespace = ? OR namespace LIKE ?)", namespace, namespace+"/%")
	var out []UnsNode
	err := q.Order("namespace").Find(&out).Error
	return out, err
}

func listUnsRecycleSubtreeByParent(db *gorm.DB, root UnsNode) ([]UnsNode, error) {
	out := make([]UnsNode, 0)
	seen := map[int64]struct{}{}
	queue := make([]int64, 0, 1)
	if root.DeletedTime > 0 && root.RecycleIsDel == 2 {
		out = append(out, root)
		seen[root.ID] = struct{}{}
		queue = append(queue, root.ID)
	} else if root.DeletedTime == 0 {
		queue = append(queue, root.ID)
	} else {
		return out, nil
	}
	for len(queue) > 0 {
		parents := queue
		queue = nil
		var children []UnsNode
		q := db.Unscoped().Where("parent_id IN ?", parents)
		if root.DeletedTime > 0 {
			q = q.Where("deleted_time > 0 AND recycle_is_del = 2")
		} else {
			q = q.Where("deleted_time = 0 OR (deleted_time > 0 AND recycle_is_del = 2)")
		}
		if err := q.
			Order("parent_id, sort_key, id").Find(&children).Error; err != nil {
			return nil, err
		}
		for _, child := range children {
			if _, ok := seen[child.ID]; ok {
				continue
			}
			seen[child.ID] = struct{}{}
			if child.DeletedTime > 0 && child.RecycleIsDel == 2 {
				out = append(out, child)
			}
			if child.Type == 1 {
				queue = append(queue, child.ID)
			}
		}
	}
	return out, nil
}

func normalizeUnsNodeDefaults(node *UnsNode) {
	if node == nil {
		return
	}
	if node.EnableHistory == 0 {
		node.EnableHistory = 2
	}
	if node.MockData == 0 {
		node.MockData = 2
	}
	if node.IsFavorite == 0 {
		node.IsFavorite = 2
	}
}

func buildUnsIDPath(db *gorm.DB, parentID, id int64) (string, error) {
	if id == 0 {
		return "", ErrInvalidUnsNode
	}
	if parentID == 0 {
		return strconv.FormatInt(id, 10), nil
	}
	parent, err := loadUnsNode(db, parentID, true)
	if err != nil {
		return "", err
	}
	parentPath := strings.Trim(parent.IDPath, "/")
	if parentPath == "" {
		parentPath, err = buildUnsIDPath(db, parent.ParentID, parent.ID)
		if err != nil {
			return "", err
		}
	}
	return parentPath + "/" + strconv.FormatInt(id, 10), nil
}

func fillUnsNodePathFields(db *gorm.DB, node *UnsNode) error {
	if node == nil {
		return ErrInvalidUnsNode
	}
	if node.ID == 0 {
		ensureID(&node.ID)
	}
	idPath, err := buildUnsIDPath(db, node.ParentID, node.ID)
	if err != nil {
		return err
	}
	node.IDPath = idPath
	return nil
}

func updateUnsIDPathPrefix(db *gorm.DB, oldPrefix, newPrefix string, userID, now int64) error {
	oldPrefix = strings.Trim(oldPrefix, "/")
	newPrefix = strings.Trim(newPrefix, "/")
	if oldPrefix == "" || newPrefix == "" || oldPrefix == newPrefix {
		return nil
	}
	return db.Exec(`
UPDATE uns_namespace_node_info
SET id_path=? || substring(id_path from ?::int), updated_by=?, updated_time=?
WHERE id_path LIKE ?`, newPrefix+"/", len(oldPrefix)+2, userID, repoTimeFromMilli(now), oldPrefix+"/%").Error
}

func recycleRestorePlan(db *gorm.DB, target UnsNode) ([]int64, error) {
	ids := make([]int64, 0)
	seen := map[int64]struct{}{}
	if target.DeletedTime > 0 {
		if target.RecycleIsDel != 2 {
			return nil, ErrUnsNotInRecycle
		}
		var ancestors []int64
		parentID := target.ParentID
		for parentID != 0 {
			parent, err := loadUnsNode(db, parentID, true)
			if err != nil {
				return nil, err
			}
			if parent.DeletedTime == 0 {
				break
			}
			if parent.RecycleIsDel != 2 {
				return nil, ErrUnsNotInRecycle
			}
			ancestors = append(ancestors, parent.ID)
			parentID = parent.ParentID
		}
		for i := len(ancestors) - 1; i >= 0; i-- {
			ids = append(ids, ancestors[i])
			seen[ancestors[i]] = struct{}{}
		}
	}
	nodes, err := listUnsRecycleSubtreeByParent(db, target)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if _, ok := seen[node.ID]; ok {
			continue
		}
		ids = append(ids, node.ID)
		seen[node.ID] = struct{}{}
	}
	if len(ids) == 0 {
		return nil, ErrUnsNotInRecycle
	}
	return ids, nil
}

func findUnsChildByName(db *gorm.DB, parentID int64, name string) (UnsNode, error) {
	var node UnsNode
	err := db.Where("parent_id = ? AND lower(name) = lower(?) AND deleted_time = 0", parentID, name).
		Order("id").Take(&node).Error
	return node, err
}

func ensureUnsSiblingNameAvailable(db *gorm.DB, parentID int64, name string, currentID int64) error {
	existing, err := findUnsChildByName(db, parentID, name)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if currentID != 0 && existing.ID == currentID {
		return nil
	}
	return ErrDuplicate
}

func buildUnsNamespace(db *gorm.DB, parentID int64, name string) (string, error) {
	name = strings.Trim(strings.TrimSpace(name), "/")
	if name == "" {
		return "", ErrInvalidUnsNode
	}
	if parentID == 0 {
		return name, nil
	}
	parent, err := loadUnsNode(db, parentID, false)
	if err != nil {
		return "", err
	}
	return strings.Trim(parent.Namespace+"/"+name, "/"), nil
}

func buildUnsNamespaceForParent(db *gorm.DB, parentID int64, name string) (string, error) {
	name = strings.Trim(strings.TrimSpace(name), "/")
	if name == "" {
		return "", ErrInvalidUnsNode
	}
	if parentID == 0 {
		return name, nil
	}
	parent, err := loadUnsNode(db, parentID, false)
	if err != nil {
		return "", err
	}
	if parent.Type != 1 {
		return "", ErrInvalidUnsNode
	}
	return strings.Trim(parent.Namespace+"/"+name, "/"), nil
}

func validateUnsFolderParent(db *gorm.DB, parentID int64) error {
	parent, err := loadUnsNode(db, parentID, false)
	if err != nil {
		return err
	}
	if parent.Type != 1 || parent.TopicType != 0 {
		return ErrInvalidUnsNode
	}
	return nil
}

func ensureUnsFileTopicParent(db *gorm.DB, node *UnsNode, now int64) error {
	enumName := unsTopicTypeName(node.TopicType)
	if enumName == "" {
		return ErrInvalidUnsNode
	}
	if node.ParentID != 0 {
		parent, err := loadUnsNode(db, node.ParentID, false)
		if err != nil {
			return err
		}
		if parent.Type != 1 {
			return ErrInvalidUnsNode
		}
		if parent.TopicType == 0 && unsTopicTypeName(node.TopicType) == parent.Name {
			node.Namespace = ""
			return nil
		}
		if parent.TopicType != 0 && parent.TopicType != node.TopicType {
			return ErrInvalidUnsNode
		}
		if parent.TopicType == node.TopicType {
			node.Namespace = ""
			return nil
		}
	}
	folder, err := findUnsChildByName(db, node.ParentID, enumName)
	if err == nil {
		if folder.Type != 1 || folder.TopicType != node.TopicType {
			return ErrInvalidUnsNode
		}
		node.ParentID = folder.ID
		node.Namespace = ""
		return nil
	}
	if !errors.Is(err, ErrNotFound) {
		return err
	}
	namespace, err := buildUnsNamespaceForParent(db, node.ParentID, enumName)
	if err != nil {
		return err
	}
	folderNode := UnsNode{
		ParentID:         node.ParentID,
		Name:             enumName,
		DisplayName:      enumName,
		Namespace:        namespace,
		Type:             1,
		TopicType:        node.TopicType,
		Schema:           json.RawMessage(fmt.Sprintf(`{"dataType":%d}`, node.TopicType)),
		ExtendProperties: json.RawMessage(`{}`),
		Status:           "active",
	}
	folderNode.CreatedBy = node.CreatedBy
	folderNode.UpdatedBy = node.CreatedBy
	ensureID(&folderNode.ID)
	normalizeUnsNodeDefaults(&folderNode)
	if err := fillUnsNodePathFields(db, &folderNode); err != nil {
		return err
	}
	if err := db.Create(&folderNode).Error; err != nil {
		return err
	}
	node.ParentID = folderNode.ID
	node.Namespace = ""
	return nil
}

func ensureRestoreFileTopicParent(db *gorm.DB, node *UnsNode, userID, now int64) error {
	enumName := unsTopicTypeName(node.TopicType)
	if enumName == "" {
		return nil
	}
	parentID := node.ParentID
	if parentID != 0 {
		parent, err := loadUnsNode(db, parentID, false)
		if err == nil {
			if parent.Type != 1 {
				return ErrInvalidUnsNode
			}
			if parent.TopicType == node.TopicType {
				return nil
			}
			if parent.TopicType != 0 {
				parentID = parent.ParentID
			}
		} else if !errors.Is(err, ErrNotFound) {
			return err
		}
	}
	folder, err := findUnsChildByName(db, parentID, enumName)
	if err == nil {
		if folder.Type != 1 || folder.TopicType != node.TopicType {
			return ErrInvalidUnsNode
		}
		node.ParentID = folder.ID
		return db.Unscoped().Model(&UnsNode{}).Where("id = ?", node.ID).
			Updates(touchByValues(map[string]any{"parent_id": folder.ID}, userID, now)).Error
	}
	if !errors.Is(err, ErrNotFound) {
		return err
	}
	namespace, err := buildUnsNamespaceForParent(db, parentID, enumName)
	if err != nil {
		return err
	}
	folderNode := UnsNode{
		ParentID:         parentID,
		Name:             enumName,
		DisplayName:      enumName,
		Namespace:        namespace,
		Type:             1,
		TopicType:        node.TopicType,
		Schema:           json.RawMessage(fmt.Sprintf(`{"dataType":%d}`, node.TopicType)),
		ExtendProperties: json.RawMessage(`{}`),
		Status:           "active",
	}
	folderNode.CreatedBy = userID
	folderNode.UpdatedBy = userID
	ensureID(&folderNode.ID)
	normalizeUnsNodeDefaults(&folderNode)
	if err := fillUnsNodePathFields(db, &folderNode); err != nil {
		return err
	}
	if err := db.Create(&folderNode).Error; err != nil {
		return err
	}
	node.ParentID = folderNode.ID
	return db.Unscoped().Model(&UnsNode{}).Where("id = ?", node.ID).
		Updates(touchByValues(map[string]any{"parent_id": folderNode.ID}, userID, now)).Error
}

func replaceUnsNodeLabels(db *gorm.DB, nodeID int64, labelIDs []int64, now int64) error {
	if err := db.Where("node_id = ?", nodeID).Delete(&unsLabelNode{}).Error; err != nil {
		return err
	}
	for _, labelID := range labelIDs {
		if labelID == 0 {
			continue
		}
		row := unsLabelNode{LabelID: labelID, NodeID: nodeID}
		row.CreatedTime = repoTimeFromMilli(now)
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "label_id"}, {Name: "node_id"}},
			DoNothing: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func bindAssetIDs(db *gorm.DB, assetIDs []int64, ownerType string, ownerID, userID, now int64) error {
	for _, assetID := range assetIDs {
		if assetID == 0 {
			continue
		}
		binding := AssetBinding{
			AssetID:   assetID,
			OwnerType: ownerType,
			OwnerID:   ownerID,
			Usage:     "attachment",
			Metadata:  json.RawMessage(`{}`),
		}
		binding.CreatedBy = userID
		binding.CreatedTime = repoTimeFromMilli(now)
		if err := db.Create(&binding).Error; err != nil {
			return err
		}
		if err := db.Model(&AssetFile{}).Where("id = ? AND deleted_time = 0", assetID).
			Updates(touchValues(map[string]any{"status": "active"}, now)).Error; err != nil {
			return err
		}
	}
	return nil
}

func availableNamespace(db *gorm.DB, base string) string {
	parent := strings.TrimSuffix(base, "/"+lastNamespacePart(base))
	name := lastNamespacePart(base)
	for i := 1; i < 1000; i++ {
		candidate := name + "-" + strconv.Itoa(i)
		if parent != "" {
			candidate = parent + "/" + candidate
		}
		var node struct {
			ID int64 `gorm:"column:id"`
		}
		err := db.Table("uns_namespace_node_info").Select("id").
			Where("namespace = ? AND deleted_time = 0", candidate).Take(&node).Error
		if errors.Is(err, ErrNotFound) {
			return candidate
		}
	}
	return base + "-" + strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func lastNamespacePart(namespace string) string {
	parts := strings.Split(strings.Trim(namespace, "/"), "/")
	if len(parts) == 0 {
		return namespace
	}
	return parts[len(parts)-1]
}

func isReservedUnsEnumFolder(node UnsNode) bool {
	return node.Type == 1 && node.TopicType != 0 && strings.EqualFold(node.Name, unsTopicTypeName(node.TopicType))
}
