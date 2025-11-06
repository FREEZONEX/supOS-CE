package dto

import "net/http"

// PaginationDTO represents pagination parameters
type PaginationDTO struct {
	PageNo   int64 `json:"pageNo" form:"pageNo"`     // 当前页数，默认为1
	PageSize int64 `json:"pageSize" form:"pageSize"` // 每页记录数，默认为20，最大支持1000
}

// GetPageNo returns page number with default value
func (p *PaginationDTO) GetPageNo() int64 {
	if p.PageNo <= 0 {
		return 1
	}
	return p.PageNo
}

// GetPageSize returns page size with default value
func (p *PaginationDTO) GetPageSize() int64 {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 1000 {
		return 1000
	}
	return p.PageSize
}

// GetOffset calculates the offset for database queries
func (p *PaginationDTO) GetOffset() int64 {
	return (p.GetPageNo() - 1) * p.GetPageSize()
}

// PageResultDTO represents paginated result
type PageResultDTO[T any] struct {
	PageNo   int64 `json:"pageNo"`        // 当前页数
	PageSize int64 `json:"pageSize"`      // 每页记录数
	Total    int64 `json:"total"`         // 总记录数
	Code     int64 `json:"code"`          // 状态码
	Data     []T   `json:"data,omitzero"` // 数据列表
}

// NewPageResult creates a new page result
func NewPageResult[T any](pageNo, pageSize, total int64, data []T) PageResultDTO[T] {
	return PageResultDTO[T]{
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Code:     http.StatusOK,
		Data:     data,
	}
}
