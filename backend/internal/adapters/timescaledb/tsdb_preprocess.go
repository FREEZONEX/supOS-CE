package timescaledb

import (
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"sort"
	"strconv"
	"time"
)

type tsData struct {
	conflictMap *base.LinkedHashMap[[2]int64, map[string]string]
	simpleMap   *base.LinkedHashMap[[2]int64, map[string]string]
}

type processResult struct {
	conflict copyParams
	normal   copyParams
}
type copyParams struct {
	columns []string
	rows    [][]any
}

// 数据预处理
func preprocess(unsData []serviceApi.UnsData) processResult {
	rs := processResult{}
	if len(unsData) == 0 {
		return rs
	}
	tsDatas := &tsData{}
	initFlags := make(map[int64]*types.CreateTopicDto, len(unsData))
	columns := make(map[string]bool, 16)
	for _, data := range unsData {
		uns, list := data.Uns, data.Data
		if len(list) == 0 || uns == nil {
			continue
		}
		def, has := initFlags[uns.Id]
		if !has {
			def = uns
			initFlags[uns.Id] = def
			uns.Timestamps[0] = -1 //先清空原来的 最大值
			if uns.Timestamps[1] == 0 {
				uns.Timestamps[1] = time.Now().UnixMilli()
			}
		}
		CT := def.GetTimestampField()
		fds := def.GetFieldDefines().FieldsMap
		for _, da := range list {
			ct := da[CT]
			mill, _ := strconv.ParseFloat(ct, 64)
			if mill > 0 {
				timeStamp := int64(mill)
				tsDatas.mergePut(def, timeStamp, da)
			}
			for k := range da {
				fd := fds[k]
				if fd != nil && fd.Index != nil {
					columns[*fd.Index] = true
				}
			}
		}
	}
	for _, tsd := range initFlags {
		if tsd.Timestamps[0] > tsd.Timestamps[1] {
			tsd.Timestamps[1] = tsd.Timestamps[0]
		}
	}
	cols := make([]string, 3, 3+len(columns))
	for k := range columns {
		cols = append(cols, k)
	}
	sort.Strings(cols[3:])
	cols[0] = constants.SysFieldCreateTime
	cols[1] = constants.SystemSeqTag
	cols[2] = constants.QosField

	colIndexMap := make(map[string]int, len(cols))
	for i, col := range cols {
		colIndexMap[col] = i
	}
	if mp := tsDatas.conflictMap; mp != nil {
		rs.conflict = copyParams{columns: cols, rows: map2rows(colIndexMap, mp, initFlags)}
	}
	if mp := tsDatas.simpleMap; mp != nil {
		rs.normal = copyParams{columns: cols, rows: map2rows(colIndexMap, mp, initFlags)}
	}
	return rs
}
func map2rows(colIndexMap map[string]int, dataMap *base.LinkedHashMap[[2]int64, map[string]string], unsMap map[int64]*types.CreateTopicDto) (rows [][]any) {
	rows = make([][]any, 0, dataMap.Size())
	dataMap.Range(func(key [2]int64, da map[string]string) {
		uns := unsMap[key[0]]
		row := make([]any, len(colIndexMap))
		tagField := uns.GetTbFieldName()
		da[tagField] = strconv.FormatInt(uns.Id, 10)
		for _, fd := range uns.Fields {
			name := fd.Name
			if index := fd.Index; index != nil && len(*index) > 0 {
				name = *index
			}
			v := da[fd.Name]
			if len(v) == 0 && fd.Type != types.FieldTypeString {
				continue
			}
			i := colIndexMap[name]
			if i >= 0 {
				if fd.Type == types.FieldTypeDatetime {
					mill, _ := strconv.ParseFloat(v, 64)
					if mill > 0 {
						utcTime := time.UnixMilli(int64(mill)).UTC()
						v = utcTime.Format("2006-01-02 15:04:05.000") + "+00"
					}
				}
				row[i] = v
			}
		}
		rows = append(rows, row)
	})
	return rows
}
func (td *tsData) mergePut(uns *types.CreateTopicDto, timestamp int64, record map[string]string) {
	//计算uns 当前时间戳最大值
	if timestamp > uns.Timestamps[0] {
		uns.Timestamps[0] = timestamp
	}
	var tsMap *base.LinkedHashMap[[2]int64, map[string]string]
	if timestamp <= uns.Timestamps[1] { //有可能主键（时间戳）冲突的
		if td.conflictMap == nil {
			td.conflictMap = base.NewLinkedHashMap[[2]int64, map[string]string]()
		}
		tsMap = td.conflictMap
	} else { // 不会冲突的
		if td.simpleMap == nil {
			td.simpleMap = base.NewLinkedHashMap[[2]int64, map[string]string]()
		}
		tsMap = td.simpleMap
	}
	key := [2]int64{uns.Id, timestamp}
	old := tsMap.Get(key)
	if len(old) > 0 {
		for k, v := range record {
			old[k] = v
		}
	} else {
		tsMap.Put(key, record)
	}
}
