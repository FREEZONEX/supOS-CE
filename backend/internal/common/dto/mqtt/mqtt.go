package mqtt

import (
	"backend/internal/types"
	"strings"
	"sync"

	"backend/internal/common/constants"
	"backend/internal/common/utils/expressionutil"
)

// TopicDefinition MQTT Topic 定义
type TopicDefinition struct {
	FieldDefines   *types.FieldDefines   `json:"fieldDefines,omitzero"`
	LastMsg        map[string]any        `json:"lastMsg,omitzero"`
	LastDt         map[string]int64      `json:"lastDt,omitzero"`
	LastDateTime   int64                 `json:"lastDateTime,omitzero"`
	ReferCalcUns   map[int64]bool        `json:"referCalcUns,omitzero"` // 被引用的计算实例（使用 map 模拟 Set）
	CreateTopicDto *types.CreateTopicDto `json:"createTopicDto,omitzero"`
	Save2db        bool                  `json:"save2db"`
	mu             sync.RWMutex          // 保护 ReferCalcUns 的并发访问
}

// NewTopicDefinition 创建新的 TopicDefinition
func NewTopicDefinition(createTopicDto *types.CreateTopicDto) *TopicDefinition {
	td := &TopicDefinition{}
	td.initByCreateTopicDto(createTopicDto, true)
	return td
}

// AddReferCalcTopic 添加引用的计算实例
func (t *TopicDefinition) AddReferCalcTopic(id int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ReferCalcUns == nil {
		t.ReferCalcUns = make(map[int64]bool, 4)
	}
	t.ReferCalcUns[id] = true
}

// RemoveReferCalcTopic 移除引用的计算实例
func (t *TopicDefinition) RemoveReferCalcTopic(id int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ReferCalcUns != nil {
		delete(t.ReferCalcUns, id)
	}
}

// GetTopic 获取 Topic
func (t *TopicDefinition) GetTopic() string {
	if t.CreateTopicDto != nil {
		return t.CreateTopicDto.GetTopic()
	}
	return ""
}

// GetTable 获取表名
func (t *TopicDefinition) GetTable() string {
	if t.CreateTopicDto != nil {
		return t.CreateTopicDto.GetTable()
	}
	return ""
}

// GetJdbcType 获取 JDBC 类型
func (t *TopicDefinition) GetJdbcType() any {
	if t.CreateTopicDto != nil {
		return t.CreateTopicDto.DataSrcID
	}
	return nil
}

// GetDataType 获取数据类型
func (t *TopicDefinition) GetDataType() int {
	if t.CreateTopicDto != nil {
		return int(*t.CreateTopicDto.DataType)
	}
	return 0
}

// GetAlarmRuleDefine 获取告警规则定义
func (t *TopicDefinition) GetAlarmRuleDefine() any {
	if t.CreateTopicDto != nil {
		return t.CreateTopicDto.AlarmRuleDefine
	}
	return nil
}

// GetRefers 获取引用字段
func (t *TopicDefinition) GetRefers() []*types.InstanceField {
	if t.CreateTopicDto != nil {
		return t.CreateTopicDto.Refers
	}
	return nil
}

// GetCompileExpression 获取编译后的表达式
func (t *TopicDefinition) GetCompileExpression() any {
	if t.CreateTopicDto == nil {
		return nil
	}

	calculationExpr := *t.CreateTopicDto.Expression
	if calculationExpr != "" && t.CreateTopicDto.CompileExpression == nil {
		compiled, _ := expressionutil.CompileExpression(calculationExpr)
		t.CreateTopicDto.CompileExpression = compiled
	}
	return t.CreateTopicDto.CompileExpression
}

// SetCreateTopicDto 设置 CreateTopicDto
func (t *TopicDefinition) SetCreateTopicDto(dto *types.CreateTopicDto) {
	t.initByCreateTopicDto(dto, false)
}

// initByCreateTopicDto 初始化
func (t *TopicDefinition) initByCreateTopicDto(createDto *types.CreateTopicDto, init bool) {
	if createDto == nil {
		return
	}

	fields := createDto.Fields
	if len(fields) > 0 {
		if createDto.GetFieldDefines() == nil {
			createDto.SetFields(fields)
		}
		fieldDefineMap := createDto.GetFieldDefines().FieldsMap
		if len(t.LastMsg) > 0 {
			for k := range t.LastMsg {
				if _, exists := fieldDefineMap[k]; !exists {
					delete(t.LastMsg, k)
					if t.LastDt != nil {
						delete(t.LastDt, k)
					}
				}
			}
		}
		t.FieldDefines = createDto.GetFieldDefines()
	} else if init {
		t.FieldDefines = types.NewFieldDefines(nil)
	} else if t.CreateTopicDto != nil {
		createDto.Fields = t.CreateTopicDto.Fields
	}

	t.CreateTopicDto = createDto

	if createDto.WithFlags != nil && *createDto.WithFlags != 0 {
		t.Save2db = constants.WithSave2db(*createDto.WithFlags)
	} else if init {
		t.Save2db = true
	}

	expr := createDto.CompileExpression
	if expr != nil || len(createDto.Refers) > 0 {
		dataType := t.GetDataType()
		if dataType == int(constants.AlarmRuleType) {
			fieldsMap := t.FieldDefines.FieldsMap
			rsField, exists := fieldsMap["isAlarm"] // AlarmRuleDefine.FIELD_IS_ALARM
			if !exists {
				// 告警表结构错误，创建默认字段
				rsField = &types.FieldDefine{
					Name: "isAlarm",
					Type: types.FieldTypeBoolean,
				}
			}
			t.FieldDefines.CalcField = rsField
		} else {
			t.FieldDefines.CalcField = getCalcField(createDto)
		}
	}
}

// getCalcField 获取计算字段
func getCalcField(createDto *types.CreateTopicDto) *types.FieldDefine {
	fields := createDto.Fields
	if len(fields) == 0 {
		return nil
	}

	var calcField *types.FieldDefine
	ct := createDto.GetTimestampField()
	qos := createDto.GetQualityField()

	for _, cv := range fields {
		name := cv.Name
		if name != ct && name != qos && !strings.HasPrefix(name, constants.SystemFieldPrev) {
			calcField = cv
			break
		}
	}
	return calcField
}
