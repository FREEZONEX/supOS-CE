package gateway

// MetaChangeReq 元数据变更请求基类
type MetaChangeReq struct {
	Size        int   `json:"size,omitzero"`        // 查询数量
	StartID     int64 `json:"startId,omitzero"`     // 起始ID
	LowVersion  int64 `json:"lowVersion,omitzero"`  // 版本下界（不含）
	HighVersion int64 `json:"highVersion,omitzero"` // 版本上界（含）(lowVersion, highVersion]
}

// CollectorMetaChangeReq 网关变更请求查询
type CollectorMetaChangeReq struct {
	MetaChangeReq
	GatewayAlias string `json:"gatewayAlias,omitzero"` // 采集器别名
	OnlyValid    bool   `json:"onlyValid,omitzero"`    // 只返回有效的数据
}

// CollectorMetaDto 采集器元数据
type CollectorMetaDto struct {
	ID          int64  `json:"id,omitzero"`          // ID
	DisplayName string `json:"displayName,omitzero"` // 显式名称
	AliasName   string `json:"aliasName,omitzero"`   // 别名（唯一）
}

// CollectorMetaChangeResp 采集器元数据变更响应
type CollectorMetaChangeResp struct {
	SaveCollectorMetaDtos   []CollectorMetaDto `json:"saveCollectorMetaDtos,omitzero"`   // 保存的采集器元数据列表
	DeleteCollectorMetaDtos []CollectorMetaDto `json:"deleteCollectorMetaDtos,omitzero"` // 删除的采集器元数据列表
}

// HasChange 检查是否有变更
func (r *CollectorMetaChangeResp) HasChange() bool {
	return len(r.SaveCollectorMetaDtos) > 0 || len(r.DeleteCollectorMetaDtos) > 0
}

// DeviceMetaChangeReq 设备（源点）变更请求
type DeviceMetaChangeReq struct {
	MetaChangeReq
	GatewayAlias string `json:"gatewayAlias,omitzero"` // 采集器别名
}

// DeviceMetaDto 设备元数据
type DeviceMetaDto struct {
	ID   int64  `json:"id,omitzero"`   // ID
	Code string `json:"code,omitzero"` // 设备编码
	Name string `json:"name,omitzero"` // 设备名称
}

// DeviceMetaChangeResp 设备（源点）变更响应
type DeviceMetaChangeResp struct {
	SaveDeviceMetaDtos   []DeviceMetaDto `json:"saveDeviceMetaDtos,omitzero"`   // 保存的设备元数据列表
	DeleteDeviceMetaDtos []DeviceMetaDto `json:"deleteDeviceMetaDtos,omitzero"` // 删除的设备元数据列表
}

// HasChange 检查是否有变更
func (r *DeviceMetaChangeResp) HasChange() bool {
	return len(r.SaveDeviceMetaDtos) > 0 || len(r.DeleteDeviceMetaDtos) > 0
}

// TagMetaChangeReq 位号元数据变更请求
type TagMetaChangeReq struct {
	MetaChangeReq
	GatewayAlias  string   `json:"gatewayAlias,omitzero"`  // 采集器别名
	DeviceAliases []string `json:"deviceAliases,omitzero"` // 设备别名集合
}

// TagMetaDto 测点位号元数据
type TagMetaDto struct {
	ID             int64  `json:"id,omitzero"`             // ID
	Code           string `json:"code,omitzero"`           // 位号编码
	Name           string `json:"name,omitzero"`           // 位号名称
	DisplayName    string `json:"displayName,omitzero"`    // 位号显示名称
	Description    string `json:"description,omitzero"`    // 位号描述
	ValueType      int    `json:"valueType,omitzero"`      // 位号值类型
	Unit           string `json:"unit,omitzero"`           // 位号单位
	Range          string `json:"range,omitzero"`          // 位号量程(下限-上限)，值域取闭区间，如0-100
	Storage        bool   `json:"storage,omitzero"`        // 是否持久化
	CollectorAlias string `json:"collectorAlias,omitzero"` // 所属采集器别名
	DeviceAlias    string `json:"deviceAlias,omitzero"`    // 所属设备别名
}

// TagMetaChangeResp 位号元数据变更响应
type TagMetaChangeResp struct {
	SaveTagMetaDtos   []TagMetaDto `json:"saveTagMetaDtos,omitzero"`   // 保存的位号元数据列表
	DeleteTagMetaDtos []TagMetaDto `json:"deleteTagMetaDtos,omitzero"` // 删除的位号元数据列表
}

// HasChange 检查是否有变更
func (r *TagMetaChangeResp) HasChange() bool {
	return len(r.SaveTagMetaDtos) > 0 || len(r.DeleteTagMetaDtos) > 0
}

// CollectorVersionReq 版本请求查询
type CollectorVersionReq struct {
	MetaChangeReq
}

// CollectorVersionResp 版本响应
type CollectorVersionResp struct {
	Version int64 `json:"version,omitzero"` // 版本号
}
