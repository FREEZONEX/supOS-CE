package dto

// CollectorGatewayQueryDto represents collector gateway query DTO
type CollectorGatewayQueryDto struct {
	PaginationDTO        // Embed PaginationDTO for pagination support
	Keyword       string `json:"keyword,omitzero" form:"keyword"` // 搜索关键字，支持显示名称和描述的模糊搜索
}
