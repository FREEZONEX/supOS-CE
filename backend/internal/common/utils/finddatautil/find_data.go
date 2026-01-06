package finddatautil

import (
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/types"
	"backend/share/base"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ListResult holds the result of a single found data list within a larger structure.
type ListResult struct {
	List       []map[string]any
	DataPath   string
	Score      int
	DataInList bool
}

// SearchResult holds the state and final outcome of the data search operation.
// It tracks the best-matching list found so far.
type SearchResult struct {
	minMatchField int
	multiFind     bool
	maxMatch      int
	List          []map[string]any
	MultiResults  []*ListResult
	DataPath      string
	DataInList    bool
	ErrorField    string
	ToLongField   string
}

// NewSearchResult creates a new SearchResult.
func NewSearchResult(minMatchField int, multiFind bool) *SearchResult {
	return &SearchResult{
		minMatchField: minMatchField,
		multiFind:     multiFind,
	}
}

// SetList updates the search result with a new potential best match.
func (sr *SearchResult) SetList(list []map[string]any, score int, dataPath string, dataInList bool) {
	if sr.multiFind {
		if sr.MultiResults == nil {
			sr.MultiResults = make([]*ListResult, 0)
		}
		sr.MultiResults = append(sr.MultiResults, &ListResult{
			List:       list,
			DataPath:   dataPath,
			Score:      score,
			DataInList: dataInList,
		})
	}

	if score > sr.maxMatch {
		sr.List = list
		sr.maxMatch = score
		sr.DataPath = dataPath
		sr.DataInList = dataInList
	}
}

// Finalize a multi-find search by sorting the results.
func (sr *SearchResult) Finalize() {
	if sr.multiFind && sr.MultiResults != nil {
		sort.Slice(sr.MultiResults, func(i, j int) bool {
			// Sort by score descending.
			return sr.MultiResults[i].Score > sr.MultiResults[j].Score
		})
	}
}

const (
	errOutOfRange = -101
)

// TypeMatchScore calculates a score based on how well the value `v` matches the `fieldType`.
// It may modify the value in `vHolder` to the converted type.
func TypeMatchScore(vHolder *any, fieldType types.FieldType, maxStrLen int) int {
	if vHolder == nil || *vHolder == nil {
		return 98
	}
	obj := *vHolder

	score := 0
	val := reflect.ValueOf(obj)

	if fieldType.IsNumber() {
		strVal := fmt.Sprintf("%v", obj)
		f, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return -1
		}
		score = 99

		switch fieldType {
		case types.FieldTypeInteger:
			if f < math.MinInt32 || f > math.MaxInt32 {
				return errOutOfRange
			}
			*vHolder = int32(f)
		case types.FieldTypeLong:
			if f < math.MinInt64 || f > math.MaxInt64 {
				return errOutOfRange
			}
			*vHolder = int64(f)
		case types.FieldTypeFloat:
			if f < -math.MaxFloat32 || f > math.MaxFloat32 {
				return errOutOfRange
			}
			*vHolder = float32(f)
		case types.FieldTypeDouble:
			*vHolder = f
		}
	} else if fieldType == types.FieldTypeString {
		s := fmt.Sprintf("%v", obj)
		if maxStrLen > 0 && len(s) > maxStrLen {
			score = errOutOfRange
		} else if val.Kind() == reflect.String {
			score = 100
		} else {
			score = 97
		}
	} else if fieldType == types.FieldTypeDatetime {
		if isNumber(val.Kind()) {
			score = 97
		} else if strVal := fmt.Sprintf("%v", obj); len(strVal) > 4 && unicode.IsDigit(rune(strVal[4])) {
			long, NaN := strconv.ParseInt(strVal, 10, 64)
			if NaN == nil {
				*vHolder = long
				score = 100
			}
		} else if timestamp, err := datetimeutils.ParseDate(strVal); err == nil {
			*vHolder = timestamp.UnixMilli() //日期统一转成时间戳毫秒
			score = 98
		}
	} else if fieldType == types.FieldTypeBoolean {
		switch val.Kind() {
		case reflect.Bool:
			score = 100
		default:
			str := strings.ToLower(fmt.Sprintf("%v", obj))
			if str == "true" {
				score = 99
				*vHolder = true
			} else if str == "false" {
				score = 99
				*vHolder = false
			} else if str == "0" {
				score = 98
				*vHolder = false
			} else if str == "1" {
				score = 98
				*vHolder = true
			} else if _, err := strconv.ParseFloat(str, 64); err == nil {
				score = 97
				*vHolder = true // Java version logic
			}
		}
	} else if fieldType == types.FieldTypeBlob || fieldType == types.FieldTypeLBlob {
		if val.Kind() == reflect.String {
			score = 100
		}
	}

	return score
}

func isSimpleType(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isNumber(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// findDataListRecursive is the main recursive function.
func findDataListRecursive(parent any, size int, obj any, fieldDefines *types.FieldDefines, rs *SearchResult, dataPath string) bool {
	if obj == nil {
		return false
	}

	val := reflect.ValueOf(obj)
	kind := val.Kind()

	if isSimpleType(kind) {
		if len(fieldDefines.FieldsMap) > 0 {
			score := -1
			var bestField string
			vHolder := obj
			for _, f := range fieldDefines.FieldsMap {
				if !f.IsSystemField() {
					tempV := vHolder
					matchScore := TypeMatchScore(&tempV, types.FieldType(f.Type), base.P2v(f.MaxLen))
					if matchScore > score {
						score = matchScore
						bestField = f.Name
						vHolder = tempV
					}
				}
			}
			if bestField != "" {
				newMap := map[string]any{bestField: vHolder}
				rs.SetList([]map[string]any{newMap}, score, "", false)
			}
		}
		return true
	}

	switch kind {
	case reflect.Slice, reflect.Array:
		if val.Len() > 0 {
			// Check if it's a slice of simple types
			firstElem := val.Index(0).Interface()
			if isSimpleType(reflect.TypeOf(firstElem).Kind()) {
				bean := make(map[string]any)
				for _, fieldDefine := range fieldDefines.FieldsMap {
					if fieldDefine.Index != nil && len(*fieldDefine.Index) > 0 {
						idx, err := strconv.Atoi(strings.TrimSpace(*fieldDefine.Index))
						if err == nil && idx < val.Len() {
							bean[fieldDefine.Name] = val.Index(idx).Interface()
						}
					}
				}
				if len(bean) > 0 {
					rs.SetList([]map[string]any{bean}, 1, dataPath, false)
				}
			} else {
				for i := 0; i < val.Len(); i++ {
					// Important: if a match is found, the loop terminates because of the 'return true'.
					if findDataListRecursive(obj, val.Len(), val.Index(i).Interface(), fieldDefines, rs, dataPath) {
						return true
					}
				}
			}
		}
	case reflect.Map:
		mapObj, ok := obj.(map[string]any)
		if !ok {
			return false
		}

		countFieldMatch := procMap(fieldDefines, rs, dataPath, mapObj)
		if countFieldMatch >= rs.minMatchField {
			var list []map[string]any
			inList := false

			// This section implements the critical logic from the Java version where if a map inside a collection
			// is found, the entire collection is filtered to keep only the maps that meet the criteria.
			if parent != nil {
				parentVal := reflect.ValueOf(parent)
				if parentVal.Kind() == reflect.Slice || parentVal.Kind() == reflect.Array {
					inList = true
					filteredList := make([]map[string]any, 0)
					for i := 0; i < parentVal.Len(); i++ {
						elem := parentVal.Index(i).Interface()
						if brotherMap, ok := elem.(map[string]any); ok {
							// Check each "brother" map in the parent collection.
							// Pass the correct dataPath for this check.
							if procMap(fieldDefines, rs, dataPath, brotherMap) >= rs.minMatchField {
								// Only add maps that meet the minimum match field criteria.
								filteredList = append(filteredList, brotherMap)
							}
						}
					}
					list = filteredList
				}
			}

			if list == nil {
				// If the parent was not a collection or was nil, the list is just the current map.
				list = []map[string]any{mapObj}
			}

			rs.SetList(list, size*2+countFieldMatch, dataPath, inList)
			return true
		}
	}

	return false
}

func procMap(fieldDefines *types.FieldDefines, rs *SearchResult, dataPath string, m map[string]any) int {
	countFieldMatch := 0
	addKvList := make([][2]any, 0)
	deleteKeys := make([]string, 0)

	for k, v := range m {
		val := reflect.ValueOf(v)
		kind := val.Kind()

		if v == nil || isSimpleType(kind) {
			vHolder := v
			var fd1, fdIdx *types.FieldDefine
			var ecA, ecB int

			fd1 = fieldDefines.FieldsMap[k]
			ecA = -2
			if fd1 != nil {
				ecA = TypeMatchScore(&vHolder, types.FieldType(fd1.Type), base.P2vWithDefault(fd1.GetMaxLen(), types.DefaultMaxStrLen))
			}

			if ecA > 0 {
				countFieldMatch++
				if vHolder != v {
					addKvList = append(addKvList, [2]any{k, vHolder})
				}
			} else {
				fieldName := fieldDefines.FieldIndexMap[k]
				if fieldName != "" {
					fdIdx = fieldDefines.FieldsMap[fieldName]
				}
				ecB = -2
				if fdIdx != nil {
					ecB = TypeMatchScore(&vHolder, types.FieldType(fdIdx.Type), base.P2vWithDefault(fdIdx.GetMaxLen(), types.DefaultMaxStrLen))
				}

				if ecB > 0 {
					countFieldMatch++
					deleteKeys = append(deleteKeys, k)
					addKvList = append(addKvList, [2]any{fieldName, vHolder})
				} else if fd1 != nil || fdIdx != nil {
					var fd *types.FieldDefine
					var errCode int
					if fd1 != nil {
						fd = fd1
						errCode = ecA
					} else {
						fd = fdIdx
						errCode = ecB
					}
					if errCode == errOutOfRange {
						rs.ToLongField = fd.Name
					} else {
						rs.ErrorField = fd.Name
					}
					return -1
				}
			}
		} else {
			var fd *types.FieldDefine
			fd = fieldDefines.FieldsMap[k]
			if fd == nil {
				fieldName := fieldDefines.FieldIndexMap[k]
				if fieldName != "" {
					fd = fieldDefines.FieldsMap[fieldName]
				}
			}

			if fd != nil && fd.Type != types.FieldTypeString {
				rs.ErrorField = fd.Name
				return -1
			}

			sc := rs.maxMatch
			newPath := k
			if dataPath != "" {
				newPath = dataPath + "." + k
			}
			findDataListRecursive(m, 0, v, fieldDefines, rs, newPath)

			if fd != nil && fd.Type == types.FieldTypeString && rs.maxMatch == sc {
				deleteKeys = append(deleteKeys, k)
				cvtStr, err := JsonUtil.ToJson(v)
				if err == nil {
					maxLen := base.P2vWithDefault(fd.GetMaxLen(), types.DefaultMaxStrLen)
					if maxLen <= 0 || len(cvtStr) <= maxLen {
						countFieldMatch++
						addKvList = append(addKvList, [2]any{fd.Name, cvtStr})
					} else {
						rs.ToLongField = fd.Name
					}
				}
			}
		}
	}

	for _, key := range deleteKeys {
		delete(m, key)
	}
	for _, kv := range addKvList {
		m[kv[0].(string)] = kv[1]
	}

	return countFieldMatch
}

// FindDataList is the main entry point for a single result search.
func FindDataList(obj any, minMatchField int, fieldDefines *types.FieldDefines) *SearchResult {
	if fieldDefines == nil {
		fieldDefines = types.NewFieldDefines(nil)
	}
	rs := NewSearchResult(minMatchField, false)
	findDataListRecursive(nil, 0, obj, fieldDefines, rs, "")
	return rs
}

// FindMultiDataList is the main entry point for finding all possible results.
func FindMultiDataList(obj any, fieldDefines *types.FieldDefines) *SearchResult {
	if fieldDefines == nil {
		fieldDefines = types.NewFieldDefines(nil)
	}
	rs := NewSearchResult(0, true)
	findDataListRecursive(nil, 0, obj, fieldDefines, rs, "")
	rs.Finalize()
	return rs
}
