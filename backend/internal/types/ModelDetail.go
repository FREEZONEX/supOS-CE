package types

// ModelDetail 实现 UnsDetail 接口
func (m *ModelDetail) SetId(id string) {
	m.Id = id
}

func (m *ModelDetail) SetTopic(topic string) {
	m.Topic = topic
}

func (m *ModelDetail) SetAlias(alias string) {
	m.Alias = alias
}

func (m *ModelDetail) SetParentAlias(parentAlias *string) {
	// ModelDetail 的 ParentAlias 是 string 类型，不是 *string
	if parentAlias != nil {
		m.ParentAlias = *parentAlias
	}
}

func (m *ModelDetail) SetPath(path string) {
	m.Path = path
}

func (m *ModelDetail) SetDataType(dataType *int16) {
	m.DataType = dataType
}

func (m *ModelDetail) SetParentDataType(parentDataType *int16) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetDataPath(dataPath *string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetPathType(pathType int16) {
	m.PathType = pathType
}

func (m *ModelDetail) SetFields(fields []*FieldDefine) {
	m.Fields = fields
}

func (m *ModelDetail) SetCreateTime(createTime int64) {
	m.CreateTime = createTime
}

func (m *ModelDetail) SetUpdateTime(updateTime int64) {
	m.UpdateTime = updateTime
}

func (m *ModelDetail) SetProtocol(protocol map[string]interface{}) {
	if m.SubscribeEnable && len(protocol) > 0 {
		if freq, has := protocol["frequency"]; has {
			m.SubscribeFrequency, _ = freq.(string)
		}
	}
}

func (m *ModelDetail) SetModelDescription(modelDescription string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetDescription(description string) {
	m.Description = description
}

func (m *ModelDetail) SetWithFlow(withFlow bool) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetWithDashboard(withDashboard bool) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetWithSave2db(withSave2db bool) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetSave2db(save2db bool) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetSubscribeEnable(subscribeEnable bool) {
	m.SubscribeEnable = subscribeEnable
}

func (m *ModelDetail) SetExpression(expression string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetShowExpression(showExpression string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetRefers(refers []InstanceField) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetLabelList(labelList []LabelVo) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetName(name string) {
	m.Name = name
}

func (m *ModelDetail) SetDisplayName(displayName string) {
	m.DisplayName = displayName
}

func (m *ModelDetail) SetPathName(pathName string) {
	m.PathName = pathName
}

func (m *ModelDetail) SetModelId(modelId string) {
	m.ModelId = modelId
}

func (m *ModelDetail) SetModelName(modelName string) {
	m.ModelName = modelName
}

func (m *ModelDetail) SetExtend(extend map[string]interface{}) {
	m.Extend = extend
}

func (m *ModelDetail) SetPayload(payload string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetTemplateName(templateName string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetTemplateAlias(templateAlias string) {
	m.TemplateAlias = templateAlias
}

func (m *ModelDetail) SetAccessLevel(accessLevel string) {
	// 空实现，ModelDetail 没有这个字段
}

func (m *ModelDetail) SetMount(mount *MountDetailVo) {
	m.Mount = mount
}
func (m *ModelDetail) SetTable(table string) {
}

func (m *ModelDetail) SetTbFieldName(name string) {
}
