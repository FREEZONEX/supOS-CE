package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/bo"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"fmt"
	"strings"
	"sync"
)

const (
	TYPE_FILE     = "file"
	TYPE_FOLDER   = "folder"
	TYPE_TEMPLATE = "template"
)
const (
	Template = "Template"
	Label    = "Label"
	UNS      = "UNS"

	Folder = "Path"
	File   = "File"
)

type FileData struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name"`
	//Namespace         string               `json:"namespace,omitempty"`
	Alias             string               `json:"alias,omitempty"`
	DisplayName       string               `json:"displayName,omitempty"`
	TemplateAlias     string               `json:"templateAlias,omitempty"`
	Fields            []*types.FieldDefine `json:"fields,omitempty"`
	DataType          string               `json:"dataType,omitempty"`
	Refers            string               `json:"refers,omitempty"`
	Expression        string               `json:"expression,omitempty"`
	Description       string               `json:"description,omitempty"`
	Label             string               `json:"label,omitempty"`
	Frequency         string               `json:"frequency,omitempty"`
	GenerateDashboard string               `json:"generateDashboard,omitempty"`
	EnableHistory     string               `json:"enableHistory,omitempty"`
	MockData          string               `json:"mockData,omitempty"`
	ParentDataType    string               `json:"topicType,omitempty"`
	Children          []*FileData          `json:"children,omitempty"`
	Error             string               `json:"error,omitempty"`

	parent *FileData
	path   string
}

func (node *FileData) getPath() string {
	if node.path == "" {
		if node.parent != nil {
			dir := node.parent.getPath()
			node.path = fmt.Sprintf("%s/%s", dir, node.Name)
		} else {
			node.path = node.Name
		}
	}
	return node.path
}
func nodeGetChildren(node *FileData) []*FileData {
	return node.Children
}
func nodeSetChildren(node *FileData, children []*FileData) {
	node.Children = children
}

func node2vo(i, parent *FileData) *types.CreateTopicDto {
	i.parent = parent
	if i.Alias == "" {
		i.Alias = PathUtil.GenerateFileAlias(i.getPath())
	}
	vo := &types.CreateTopicDto{
		Alias:  i.Alias,
		Name:   i.Name,
		Fields: i.Fields,
		Path:   i.getPath(),
	}

	if parent != nil {
		if parent.Alias == "" {
			parent.Alias = PathUtil.GenerateFileAlias(parent.getPath())
		}
		vo.ParentAlias = &parent.Alias
	}
	if len(i.DisplayName) > 0 {
		vo.DisplayName = &i.DisplayName
	}
	if len(i.TemplateAlias) > 0 {
		vo.ModelAlias = &i.TemplateAlias
	}
	if len(i.Description) > 0 {
		vo.Description = &i.Description
	}
	if len(i.Label) > 0 {
		vo.LabelNames = strings.Split(i.Label, ",")
	}
	switch i.Type {
	case TYPE_FILE:
		vo.PathType = constants.PathTypeFile
	case TYPE_FOLDER:
		vo.PathType = constants.PathTypeDir
	case TYPE_TEMPLATE:
		vo.PathType = constants.PathTypeTemplate
	default:
		vo.PathType = -1
	}

	if len(i.ParentDataType) > 0 {
		if pdt, ok := enums.GetFolderDataTypeByName(i.ParentDataType); ok {
			vo.ParentDataType = base.V2p(int16(pdt))
		}
	}
	if len(i.DataType) > 0 {
		dt := enums.DataTypeInt(i.DataType)
		if dt >= 0 {
			vo.DataType = base.V2p(dt)
		}
	}
	return vo
}
func poGetId(node *dao.UnsNamespace) int64 {
	return node.Id
}
func poGetParentId(node *dao.UnsNamespace) int64 {
	return base.P2vWithDefault(node.ParentId, -1)
}

var initDefOnce sync.Once
var defService serviceApi.IUnsDefinitionService

func po2DataVo(unsPo *dao.UnsNamespace) *FileData {
	return uns2DataVo(unsPo)
}
func vo2DataVo(unsPo *types.CreateTopicDto) *FileData {
	return uns2DataVo(unsPo)
}
func uns2DataVo(unsPo bo.UnsInfo) *FileData {

	data := &FileData{}

	data.Alias = unsPo.GetAlias()
	data.DisplayName = unsPo.GetDisplayName()
	if mid := unsPo.GetModelId(); mid != nil {
		if defService == nil {
			initDefOnce.Do(func() {
				defService = spring.GetBean[serviceApi.IUnsDefinitionService]()
			})
		}
		template := defService.GetDefinitionById(*mid)
		if template != nil {
			data.TemplateAlias = template.Alias
		}
	}
	// data.Namespace = unsPo.Path
	data.Name = unsPo.GetName()
	data.Expression = unsPo.GetExpression()

	if labels := unsPo.GetLabelIds(); len(labels) > 0 {
		data.Label = strings.Join(base.MapValues(labels), ",")
	}

	if protocol := unsPo.GetProtocolMap(); len(protocol) > 0 {
		frequency := protocol["frequency"]
		if frequency != nil {
			data.Frequency = fmt.Sprint(frequency)
		}
	}

	data.Description = unsPo.GetDescription()
	if pdt := unsPo.GetParentDataType(); pdt != nil {
		data.ParentDataType = enums.GetFolderDataType(*pdt).Name()
	}
	if dt := unsPo.GetDataType(); dt != nil {
		pt := unsPo.GetPathType()
		if pt == constants.PathTypeFile {
			data.DataType = enums.DataTypeName(*dt)
		}
		if *dt != constants.TimeSequenceType {
			data.Fields = base.Filter(unsPo.GetFields(), func(e *types.FieldDefine) bool {
				return !e.IsSystemField()
			})
		} else {
			data.Fields = unsPo.GetFields()
		}
	}

	//if unsPo.DataType == constants.CALCULATION_REAL_TYPE ||
	//	unsPo.DataType == constants.MERGE_TYPE ||
	//	unsPo.DataType == constants.CITING_TYPE {
	//	data.Refers = handleRefer(context, unsPo.Refers, unsPo.DataType)
	//}

	//if unsPo.DataType == constants.MERGE_TYPE || unsPo.DataType == constants.CITING_TYPE {
	//	data.Fields = nil
	//}
	switch unsPo.GetPathType() {
	case constants.PathTypeFile:
		data.Type = TYPE_FILE
		flags := unsPo.GetFlags()
		if flags != nil {
			fl := *flags
			data.EnableHistory = _BOOL(constants.WithSave2db(fl))
			data.GenerateDashboard = _BOOL(constants.WithDashBoard(fl))
			data.MockData = _BOOL(constants.WithFlow(fl))
		}
	case constants.PathTypeDir:
		data.Type = TYPE_FOLDER
	}

	return data
}
func _BOOL(b bool) string {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
}
