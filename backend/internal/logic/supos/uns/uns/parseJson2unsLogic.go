// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/common/utils/finddatautil"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"encoding/json"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseJson2unsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部JSON定义转uns字段定义
func NewParseJson2unsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseJson2unsLogic {
	return &ParseJson2unsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseJson2unsLogic) ParseJson2uns(req []byte) (resp *types.ParseJson2UnsResp, err error) {
	var data interface{}
	err = json.Unmarshal(req, &data)
	resp = &types.ParseJson2UnsResp{}
	resp.Code, resp.Msg = 200, "OK"
	if err != nil {
		resp.Code, resp.Msg = 400, "NotJson: "+err.Error()
		return
	}
	searchResult := finddatautil.FindMultiDataList(data, nil)
	if searchResult == nil || len(searchResult.MultiResults) == 0 {
		resp.Code, resp.Msg = 404, "NotFindDataList"
		return
	}
	onlyLists := make([]*finddatautil.ListResult, 0, len(searchResult.MultiResults))
	for _, result := range searchResult.MultiResults {
		if result.DataInList {
			onlyLists = append(onlyLists, result)
		}
	}
	list := onlyLists
	if len(onlyLists) == 0 {
		list = searchResult.MultiResults
	}
	pathMap := make(map[string]*types.OuterStructureVo, len(list))
	for _, rs := range list {
		pathMap[rs.DataPath] = map2fields(rs)
	}
	paths := base.MapKeys(pathMap)
	sort.Strings(paths)
	rsList := make([]*types.OuterStructureVo, 0, len(paths))
	for _, path := range paths {
		rsList = append(rsList, pathMap[path])
	}
	resp.Data = rsList
	return
}

// map2fields 函数
func map2fields(rs *finddatautil.ListResult) *types.OuterStructureVo {
	fieldTypes := make(map[string]types.FieldType)
	values := make(map[string]interface{})
	fieldOrder := make([]string, 0) // 用于保持字段顺序

	for _, data := range rs.List {
		for k, v := range data {
			fieldType, exists := fieldTypes[k]
			_guessType := guessType(v)

			if !exists || _guessType.GetOrdinal() > fieldType.GetOrdinal() {
				fieldTypes[k] = _guessType
				values[k] = v
				// 如果是新字段，添加到顺序列表中
				if !exists {
					fieldOrder = append(fieldOrder, k)
				}
			}
		}
	}

	if len(fieldTypes) == 0 {
		return nil
	}

	// 按照原始顺序构建字段列表
	fields := make([]*types.FieldDefine, 0, len(fieldOrder))
	for _, k := range fieldOrder {
		fields = append(fields, &types.FieldDefine{
			Name: k,
			Type: fieldTypes[k].Name(),
		})
	}

	return &types.OuterStructureVo{
		DataPath: rs.DataPath,
		Fields:   fields,
		ValueMap: values,
	}
}

// guessType 函数
func guessType(o interface{}) types.FieldType {
	if o == nil {
		return types.FieldTypeString
	}

	switch v := o.(type) {
	case int:
		return types.FieldTypeInteger
	case int32:
		return types.FieldTypeInteger
	case int64:
		return types.FieldTypeLong
	case float32:
		return types.FieldTypeFloat
	case float64:
		return types.FieldTypeDouble
	case *big.Float:
		return types.FieldTypeDouble
	case *big.Rat:
		return types.FieldTypeDouble
	case bool:
		return types.FieldTypeBoolean
	case string:
		_, er := strconv.ParseInt(v, 10, 64)
		if er == nil {
			return types.FieldTypeLong
		}
		_, er = strconv.ParseFloat(v, 64)
		if er == nil {
			return types.FieldTypeDouble
		}
		if parseDate(v) != nil {
			return types.FieldTypeDatetime
		}
		return types.FieldTypeString
	default:
		// 尝试将其他类型转换为字符串再判断
		str := toString(v)
		if parseDate(str) != nil {
			return types.FieldTypeDatetime
		}
		return types.FieldTypeString
	}
}

// parseDate 尝试解析日期字符串
func parseDate(dateStr string) *time.Time {
	// 常见的日期格式
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t
		}
	}

	// 尝试 Unix 时间戳
	if timestamp, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
		t := time.Unix(timestamp, 0)
		return &t
	}

	return nil
}

// toString 将任意值转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int32, int64, float32, float64, bool:
		return strings.TrimSpace(strings.Trim(strings.ReplaceAll(strconv.Quote(strings.TrimSpace(strconv.FormatFloat(float64(0), 'f', -1, 64))), "\"", ""), " "))
	default:
		return ""
	}
}
