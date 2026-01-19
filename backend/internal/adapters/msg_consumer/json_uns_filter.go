package msg_consumer

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/types"
	"backend/share/base"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

func filterMsgByUns(def *types.UnsDefinition, dataList []map[string]string) (errMsg string) {
	fds := def.GetFieldDefines().FieldsMap
	CT := def.GetTimestampField()
	onlyNormalField, countNormals := "", 0
	for i, dm := range dataList {
		for k, v := range dm {
			fd, has := fds[k]
			if !has {
				if len(dm) == 1 {
					if onlyNormalField == "" {
						for _, fieldDef := range fds {
							if !fieldDef.IsSystemField() {
								onlyNormalField = fieldDef.Name
								countNormals++
							}
						}
					}
					if countNormals == 1 {
						delete(dm, k)
						k = onlyNormalField
						dm[k] = v
					} else {
						continue
					}
				} else {
					continue
				}
			}
			fieldType := types.FieldType(fd.Type)
			newValue, errTip, qos := checkFieldValue(fieldType, v, base.P2vWithDefault(fd.MaxLen, -1))
			if errTip != "" {
				errMsg += errTip + " "
				if qos != 0 {
					qosField := def.GetQualityField()
					var objMap = make(map[string]string, 4)
					objMap[CT] = dm[CT]
					var defVal = fd.GetType().DefaultValueStr()
					objMap[k] = defVal
					objMap[qosField] = strconv.FormatInt(qos, 10)
					dataList[i] = objMap
					break
				}
			}
			if newValue != "" {
				dm[k] = newValue
			}
		}
	}
	return
}
func checkFieldValue(fieldType types.FieldType, v string, maxStrLen int) (newValue, errMsg string, qos int64) {
	errField, errOutOfRange := "", ""
	if fieldType.IsNumber() {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errField = err.Error()
		} else {
			switch fieldType {
			case types.FieldTypeInteger:
				if f < math.MinInt32 || f > math.MaxInt32 {
					errOutOfRange = "1"
				}
			case types.FieldTypeLong:
				if f < math.MinInt64 || f > math.MaxInt64 {
					errOutOfRange = "1"
				}
			case types.FieldTypeFloat:
				if f < -math.MaxFloat32 || f > math.MaxFloat32 {
					errOutOfRange = "1"
				}
			case types.FieldTypeDouble:
				if f < -math.MaxFloat64 || f > math.MaxFloat64 {
					errOutOfRange = "1"
				}
			}
		}
	} else if fieldType == types.FieldTypeString {
		if maxStrLen > 0 && len(v) > maxStrLen {
			errOutOfRange = "1"
		}
	} else if fieldType == types.FieldTypeDatetime {
		if len(v) > 4 && unicode.IsDigit(rune(v[4])) {
			Float, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errField = err.Error()
			} else if ct := int64(Float); ct < 1100000000000 || ct > 11000000000001 {
				errOutOfRange = "1"
			} else if fastInt64Len(ct) != len(v) {
				newValue = strconv.FormatInt(ct, 10)
			}
		} else if timestamp, err := datetimeutils.ParseDate(v); err == nil {
			//日期统一转成时间戳毫秒
			newValue = strconv.FormatInt(timestamp.UnixMilli(), 10)
		} else {
			errField = "Not A DateTime"
		}
	} else if fieldType == types.FieldTypeBoolean {
		if v != "true" && v != "false" {
			Bool, err := strconv.ParseBool(v)
			if err != nil {
				errField = err.Error()
			} else {
				newValue = strconv.FormatBool(Bool)
			}
		}
	}
	if len(errField) > 0 {
		qos = 0x400000000000000
		errMsg = I18nUtils.GetMessage("uns.invalid.type", v)
	}
	if len(errOutOfRange) > 0 {
		qos = 0x80000000000000 //超量程（工程单位）值"
		var tip = fmt.Sprintf("%s...len=%d", v[:min(10, len(v))], len(v))
		errMsg = I18nUtils.GetMessage("uns.invalid.toLong", tip)
	}
	return
}

var digitBounds = [19]int64{
	9, 99, 999, 9999, 99999, 999999, 9999999, 99999999, 999999999,
	9999999999, 99999999999, 999999999999, 9999999999999,
	99999999999999, 999999999999999, 9999999999999999,
	99999999999999999, 999999999999999999, 9223372036854775807,
}

func fastInt64Len(i int64) int {
	if i >= 0 {
		return positiveLen(i)
	}
	// 负数：长度 = 正数长度 + 1（负号）
	if i == -9223372036854775808 {
		return 20 // 特殊处理最小值
	}
	return positiveLen(-i) + 1
}

func positiveLen(n int64) int {
	if n == 0 {
		return 1
	}

	// 二分查找
	left, right := 0, len(digitBounds)
	for left < right {
		mid := (left + right) / 2
		if n <= digitBounds[mid] {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left + 1
}
