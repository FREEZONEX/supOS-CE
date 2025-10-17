package relationDB

import (
	"backend/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (p UnsNamespaceRepo) UpdateModelFieldsById(ctx context.Context, id uint, fields []types.FieldDefine, numberCount int, updateAt time.Time) (int64, error) {
	db := p.dbx(ctx)
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return 0, fmt.Errorf("marshal fields error: %v", err)
	}
	result := db.
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"table_name":    TableNameUnsNamespace,
			"fields":        fieldsJSON,
			"number_fields": numberCount,
			"update_at":     updateAt,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
