package dto

// SysModuleDto represents system module DTO
type SysModuleDto struct {
	ModuleCode string `json:"moduleCode,omitzero"` // 模块编码
	ModuleName string `json:"moduleName,omitzero"` // 模块名称
}
