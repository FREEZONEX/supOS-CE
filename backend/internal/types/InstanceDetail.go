package types

import "backend/internal/common/utils/JsonUtil"

// InstanceDetail 实现 UnsDetail 接口
func (i *InstanceDetail) SetId(id string) {
	i.Id = id
}

func (i *InstanceDetail) SetTopic(topic string) {
	i.Topic = topic
}

func (i *InstanceDetail) SetAlias(alias string) {
	i.Alias = alias
}

func (i *InstanceDetail) SetParentAlias(parentAlias *string) {
	i.ParentAlias = parentAlias
}

func (i *InstanceDetail) SetPath(path string) {
	i.Path = path
}

func (i *InstanceDetail) SetDataType(dataType *int16) {
	i.DataType = dataType
}

func (i *InstanceDetail) SetParentDataType(parentDataType *int16) {
	i.ParentDataType = parentDataType
}

func (i *InstanceDetail) SetDataPath(dataPath *string) {
	i.DataPath = dataPath
}

func (i *InstanceDetail) SetPathType(pathType int16) {
	i.PathType = pathType
}

func (i *InstanceDetail) SetFields(fields []*FieldDefine) {
	i.Fields = fields
}

func (i *InstanceDetail) SetCreateTime(createTime int64) {
	i.CreateTime = createTime
}

func (i *InstanceDetail) SetUpdateTime(updateTime int64) {
	i.UpdateTime = updateTime
}

func (i *InstanceDetail) SetProtocol(protocol map[string]interface{}) {
	i.Protocol = protocol
	if len(protocol) > 0 {
		if jsf, has := protocol["jsf"]; has {
			delete(protocol, "jsf")
			jsonFs := ""
			if str, isStr := jsf.(string); isStr {
				jsonFs = str
			} else {
				jsonFs, _ = JsonUtil.ToJson(jsf)
			}
			if len(jsonFs) > 0 {
				_ = JsonUtil.FromJson(jsonFs, &i.JsonFields)
			}
		}
	}
}

func (i *InstanceDetail) SetModelDescription(modelDescription string) {
	i.ModelDescription = modelDescription
}

func (i *InstanceDetail) SetDescription(description string) {
	i.Description = description
}

func (i *InstanceDetail) SetWithFlow(withFlow bool) {
	i.WithFlow = withFlow
}

func (i *InstanceDetail) SetWithDashboard(withDashboard bool) {
	i.WithDashboard = withDashboard
}

func (i *InstanceDetail) SetWithSave2db(withSave2db bool) {
	i.WithSave2db = withSave2db
}

func (i *InstanceDetail) SetSave2db(save2db bool) {
	i.Save2db = save2db
}

func (i *InstanceDetail) SetSubscribeEnable(subscribeEnable bool) {
	i.SubscribeEnable = subscribeEnable
}

func (i *InstanceDetail) SetExpression(expression string) {
	i.Expression = expression
}

func (i *InstanceDetail) SetShowExpression(showExpression string) {
	i.ShowExpression = showExpression
}

func (i *InstanceDetail) SetRefers(refers []InstanceField) {
	i.Refers = refers
}

func (i *InstanceDetail) SetLabelList(labelList []LabelVo) {
	i.LabelList = labelList
}

func (i *InstanceDetail) SetName(name string) {
	i.Name = name
}

func (i *InstanceDetail) SetDisplayName(displayName string) {
	i.DisplayName = displayName
}

func (i *InstanceDetail) SetPathName(pathName string) {
	i.PathName = pathName
}

func (i *InstanceDetail) SetModelId(modelId string) {
	i.ModelId = modelId
}

func (i *InstanceDetail) SetModelName(modelName string) {
	i.ModelName = modelName
}

func (i *InstanceDetail) SetExtend(extend map[string]interface{}) {
	i.Extend = extend
}

func (i *InstanceDetail) SetPayload(payload string) {
	i.Payload = payload
}

func (i *InstanceDetail) SetTemplateName(templateName string) {
	i.TemplateName = templateName
}

func (i *InstanceDetail) SetTemplateAlias(templateAlias string) {
	i.TemplateAlias = templateAlias
}

func (i *InstanceDetail) SetAccessLevel(accessLevel string) {
	i.AccessLevel = accessLevel
}

func (i *InstanceDetail) SetMount(mount *MountDetailVo) {
	i.Mount = mount
}
func (i *InstanceDetail) SetTable(table string) {
	i.Table = table
}

func (i *InstanceDetail) SetTbFieldName(name string) {
	i.TbFieldName = name
}
