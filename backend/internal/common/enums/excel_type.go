package enums

import (
	"backend/internal/common/constants"
	"sort"
)

// ExcelTypeEnum represents Excel sheet types for import/export
type ExcelTypeEnum struct {
	Code  string
	Index int
}

var (
	ExcelTypeExplanation = ExcelTypeEnum{Code: "Explanation", Index: 0}
	ExcelTypeTemplate    = ExcelTypeEnum{Code: "Template", Index: 1}
	ExcelTypeLabel       = ExcelTypeEnum{Code: "Label", Index: 2}
	ExcelTypeFolder      = ExcelTypeEnum{Code: "Path", Index: 3}

	ExcelTypeFileTimeseries  = ExcelTypeEnum{Code: "Topic-timeseries", Index: 4}
	ExcelTypeFileRelation    = ExcelTypeEnum{Code: "Topic-relation", Index: 5}
	ExcelTypeFileCalculate   = ExcelTypeEnum{Code: "Topic-calculate", Index: 6}
	ExcelTypeFileAggregation = ExcelTypeEnum{Code: "Topic-aggregation", Index: 7}
	ExcelTypeFileReference   = ExcelTypeEnum{Code: "Topic-reference", Index: 8}

	ExcelTypeUNS   = ExcelTypeEnum{Code: "UNS", Index: 3}
	ExcelTypeFile  = ExcelTypeEnum{Code: "File", Index: 4}
	ExcelTypeError = ExcelTypeEnum{Code: "error", Index: -1}
)

// allExcelTypes contains all Excel types
var allExcelTypes = []ExcelTypeEnum{
	ExcelTypeExplanation,
	ExcelTypeTemplate,
	ExcelTypeLabel,
	ExcelTypeFolder,
	ExcelTypeFileTimeseries,
	ExcelTypeFileRelation,
	ExcelTypeFileCalculate,
	ExcelTypeFileAggregation,
	ExcelTypeFileReference,
	ExcelTypeUNS,
	ExcelTypeFile,
	ExcelTypeError,
}

// GetExcelTypeFromCode returns ExcelTypeEnum from code
func GetExcelTypeFromCode(code string) ExcelTypeEnum {
	for _, et := range allExcelTypes {
		if et.Code == code {
			return et
		}
	}
	return ExcelTypeError
}

// GetExcelTypeFromIndex returns ExcelTypeEnum from index
func GetExcelTypeFromIndex(index int) ExcelTypeEnum {
	for _, et := range allExcelTypes {
		if et.Index == index {
			return et
		}
	}
	return ExcelTypeError
}

// ExcelTypeSorted returns sorted Excel types by index
func ExcelTypeSorted() []ExcelTypeEnum {
	sorted := make([]ExcelTypeEnum, len(allExcelTypes))
	copy(sorted, allExcelTypes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Index < sorted[j].Index
	})

	return sorted
}

// ExcelTypeListFile returns file-type Excel types
func ExcelTypeListFile() []ExcelTypeEnum {
	return []ExcelTypeEnum{
		ExcelTypeFileTimeseries,
		ExcelTypeFileRelation,
		ExcelTypeFileCalculate,
		ExcelTypeFileAggregation,
		ExcelTypeFileReference,
	}
}

// GetExcelTypeFromDataType returns ExcelTypeEnum from data type constant
func GetExcelTypeFromDataType(dataType int16) ExcelTypeEnum {
	switch dataType {
	case constants.TimeSequenceType:
		return ExcelTypeFileTimeseries
	case constants.RelationType:
		return ExcelTypeFileRelation
	case constants.CalculationRealType:
		return ExcelTypeFileCalculate
	case constants.MergeType:
		return ExcelTypeFileAggregation
	case constants.CitingType:
		return ExcelTypeFileReference
	default:
		return ExcelTypeError
	}
}

// ExcelTypeSize returns the count of Excel types (excluding error type)
func ExcelTypeSize() int {
	return len(allExcelTypes) - 1
}
