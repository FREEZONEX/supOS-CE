package bo

import (
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
)

// UnsDetail 接口定义
type UnsDetail interface {
	// 基础字段
	SetId(id string)
	SetTopic(topic string)
	SetAlias(alias string)
	SetParentAlias(parentAlias *string)
	SetPath(path string)
	SetDataType(dataType *int16)
	SetParentDataType(parentDataType *int16)
	SetDataPath(dataPath *string)
	SetPathType(pathType int16)
	SetFields(fields []*types.FieldDefine)
	SetCreateTime(createTime int64)
	SetUpdateTime(updateTime int64)
	SetProtocol(protocol map[string]interface{})
	SetModelDescription(modelDescription string)
	SetDescription(description string)
	SetWithFlow(withFlow bool)
	SetWithDashboard(withDashboard bool)
	SetWithSave2db(withSave2db bool)
	SetSave2db(save2db bool)
	SetSubscribeEnable(subscribeEnable bool)
	SetExpression(expression string)
	SetShowExpression(showExpression string)
	SetRefers(refers []types.InstanceField)
	SetLabelList(labelList []types.LabelVo)
	SetName(name string)
	SetDisplayName(displayName string)
	SetPathName(pathName string)
	SetModelId(modelId string)
	SetModelName(modelName string)
	SetTable(table string)
	SetTbFieldName(name string)
	SetExtend(extend map[string]interface{})
	SetPayload(payload string)
	SetTemplateName(templateName string)
	SetTemplateAlias(templateAlias string)
	SetAccessLevel(accessLevel string)
	SetMount(mount *types.MountDetailVo)
}

var _ UnsDetail = &types.InstanceDetail{}
var _ UnsDetail = &types.ModelDetail{}
var _ types.UnsInfo = &types.CreateTopicDto{}
var _ types.UnsInfo = &dao.UnsNamespace{}
