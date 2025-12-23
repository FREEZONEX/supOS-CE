package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/PathUtil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	TYPE_FILE     = "topic"
	TYPE_FOLDER   = "path"
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
	JsonFields        []*types.FieldDefine `json:"jsonFields,omitempty"`
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

	parent   *FileData
	path     string
	id       int64
	parentId int64
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
func nodeGetId(node *FileData) int64 {
	return node.id
}
func nodeGetParentId(node *FileData) int64 {
	return node.parentId
}

func node2vo(prop string, i, parent *FileData) *types.CreateTopicDto {
	if i.Name == "" {
		i.Error = "Empty " + prop
		return nil
	}
	if prop == Label {
		return &types.CreateTopicDto{Name: i.Name}
	}

	vo := &types.CreateTopicDto{
		Alias:      i.Alias,
		Name:       i.Name,
		Fields:     i.Fields,
		JsonFields: i.JsonFields,
	}
	switch prop {
	case Template:
		vo.PathType = constants.PathTypeTemplate
	case UNS:
		switch i.Type {
		case TYPE_FILE, "file":
			vo.PathType = constants.PathTypeFile
			if os.Getenv("SYS_OS_ENABLE_AUTO_CATEGORIZATION") == "true" {
				if i.ParentDataType == "" {
					i.Error = I18nUtils.GetMessage("uns.excel.parentDataType.is.blank")
					return nil
				}
			}
		case TYPE_FOLDER, "folder":
			vo.PathType = constants.PathTypeDir
			if i.Name == "label" || i.Name == "template" {
				i.Error = I18nUtils.GetMessage("uns.folder.reserved.word")
				return nil
			}
			if len(i.Name) > 63 {
				i.Error = I18nUtils.GetMessage("uns.folder.length.limit.exceed")
				return nil
			}
		default:
			i.Error = I18nUtils.GetMessage("uns.import.type.error")
			return nil
		}
	}

	if vo.Alias == "" {
		vo.Alias = PathUtil.GenerateFileAlias(i.getPath())
	}
	i.parent = parent
	if parent != nil {
		if parent.Alias == "" {
			parent.Alias = PathUtil.GenerateFileAlias(parent.getPath())
		}
		vo.ParentAlias = &parent.Alias
	}
	vo.Path = i.getPath()
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
	if len(i.DataType) > 0 {
		dt := enums.DataTypeInt(i.DataType)
		if dt >= 0 {
			vo.DataType = base.V2p(dt)
		} else {
			i.Error = I18nUtils.GetMessage("uns.import.dataType.error")
			return nil
		}
	}
	if len(i.ParentDataType) > 0 {
		if pdt, ok := enums.GetFolderDataTypeByName(i.ParentDataType); ok {
			dirType := base.V2p(int16(pdt))
			switch vo.PathType {
			case constants.PathTypeDir:
				vo.DataType = dirType
			case constants.PathTypeFile:
				vo.ParentDataType = dirType
			}
		}
	}
	return vo
}

var initDefOnce sync.Once
var defService serviceApi.IUnsDefinitionService

func po2DataVo(unsPo *dao.UnsNamespace) *FileData {
	return uns2DataVo(unsPo)
}
func vo2DataVo(unsPo *types.CreateTopicDto) *FileData {
	return uns2DataVo(unsPo)
}
func uns2DataVo(unsPo types.UnsInfo) *FileData {

	data := &FileData{id: unsPo.GetId(), parentId: base.P2vWithDefault(unsPo.GetParentId(), -1)}

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
