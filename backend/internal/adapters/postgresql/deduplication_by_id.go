package postgresql

import (
	"backend/internal/types"
	"backend/share/base"
	"fmt"
)

func DeduplicationById(def *types.CreateTopicDto, data []map[string]string) []map[string]string {
	pks := def.GetPrimaryField()
	if len(pks) == 0 {
		return data
	}
	idMap := base.NewLinkedHashMap[string, map[string]string]()
	mergePut := func(id string, record map[string]string) {
		old := idMap.Get(id)
		if len(old) > 0 {
			for k, v := range record {
				old[k] = v
			}
		} else {
			idMap.Put(id, record)
		}
	}
	if sz := len(pks); sz == 1 {
		pk := pks[0]
		for _, record := range data {
			var id string
			if val, ok := record[pk]; ok {
				id = val
			} else {
				id = fmt.Sprintf("%p", record)
			}
			mergePut(id, record)
		}
	} else {
		initCap := sz * 20
		var idBd base.StringBuilder
		idBd.Grow(initCap)
		for _, record := range data {
			allEmpty := true
			for _, pk := range pks {
				if val, ok := record[pk]; ok {
					allEmpty = false
					idBd.Append(val).Append("`")
				} else {
					idBd.Append("null`")
				}
			}
			var id string
			if allEmpty {
				id = fmt.Sprintf("%p", record)
			} else {
				id = idBd.String()
			}
			idBd.Reset()
			mergePut(id, record)
		}
	}
	return idMap.Values()
}
