package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"fmt"
	"os"
	"strconv"
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
	Template          *FileData            `json:"template,omitempty"`
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
		Alias:  i.Alias,
		Name:   i.Name,
		Fields: i.Fields,
	}
	switch prop {
	case Template:
		vo.PathType = constants.PathTypeTemplate
	case UNS:
		switch strings.ToLower(i.Type) {
		case TYPE_FILE:
			vo.PathType = constants.PathTypeFile
			if ok, _ := strconv.ParseBool(os.Getenv("SYS_OS_ENABLE_AUTO_CATEGORIZATION")); ok {
				if i.ParentDataType == "" {
					i.Error = I18nUtils.GetMessage("uns.excel.parentDataType.is.blank")
					return nil
				}
			}
		case TYPE_FOLDER:
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
	if template := i.Template; template != nil {
		vo.Template = node2vo(Template, template, nil)
	}
	if len(i.Description) > 0 {
		vo.Description = &i.Description
	}
	if len(i.Label) > 0 {
		vo.LabelNames = base.FilterAndMap(strings.Split(i.Label, ","), func(s string) (string, bool) {
			s = strings.TrimSpace(s)
			return s, len(s) > 0
		})
	}
	if dt := i.DataType; len(dt) > 0 {
		switch vo.PathType {
		case constants.PathTypeFile:
			dt := enums.DataTypeInt(dt)
			if dt >= 0 {
				vo.DataType = base.V2p(dt)
				if dt == constants.JsonbType {
					vo.Fields = nil
					vo.JsonFields = i.Fields
				}
			} else {
				i.Error = I18nUtils.GetMessage("uns.import.dataType.error")
				return nil
			}
		case constants.PathTypeDir:
			if pdt, ok := enums.GetFolderDataTypeByName(dt); ok {
				vo.DataType = base.V2p(pdt.TypeIndex())
			}
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
	if vo.PathType == constants.PathTypeFile {
		vo.AddDashBoard = parseBoolP(i.GenerateDashboard)
		vo.Save2Db = parseBoolP(i.EnableHistory)

		if base.P2v(vo.DataType) == constants.JsonbType {
			vo.AddFlow = base.OptionalFalse
		} else {
			vo.AddFlow = parseBoolP(i.MockData)
		}
	}
	return vo
}
func parseBoolP(str string) *bool {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	if v, er := strconv.ParseBool(str); er == nil {
		return &v
	} else {
		return base.OptionalFalse
	}
}

var initDefOnce sync.Once
var defService serviceApi.IUnsDefinitionService

type exportContext struct {
	unsWrapTemplate bool // UNS 文件内部是否嵌套模版的定义
	templates       map[int64]*FileData
}

func uns2DataVo(ctx *exportContext, unsPo types.UnsInfo) *FileData {

	data := &FileData{id: unsPo.GetId(), parentId: base.P2vWithDefault(unsPo.GetParentId(), -1)}

	data.Alias = unsPo.GetAlias()
	data.DisplayName = unsPo.GetDisplayName()
	if mid := unsPo.GetModelId(); mid != nil && ctx != nil {
		if defService == nil {
			initDefOnce.Do(func() {
				defService = spring.GetBean[serviceApi.IUnsDefinitionService]()
			})
		}
		if templateId := *mid; ctx.unsWrapTemplate {
			if ctx.templates == nil {
				ctx.templates = make(map[int64]*FileData, 16)
			}
			if v, has := ctx.templates[templateId]; has {
				data.TemplateAlias = v.Alias
			} else {
				if template := defService.GetDefinitionById(templateId); template != nil {
					templateVo := uns2DataVo(nil, template)
					data.Template = templateVo //首个引用该模版的uns, 嵌套模版的定义
					ctx.templates[templateId] = templateVo
				} else {
					ctx.templates[templateId] = nil
				}
			}
		} else {
			if template := defService.GetDefinitionById(templateId); template != nil {
				data.TemplateAlias = template.Alias
			}
		}
	}
	// data.Namespace = unsPo.Path
	data.Name = unsPo.GetName()
	data.Expression = unsPo.GetExpression()

	if labels := unsPo.GetLabelIds(); len(labels) > 0 {
		data.Label = strings.Join(base.MapValues(labels), ",")
	}

	if dt := unsPo.GetDataType(); dt != nil {
		pt := unsPo.GetPathType()
		if pt == constants.PathTypeFile {
			data.DataType = enums.DataTypeName(*dt)
		} else if pt == constants.PathTypeDir {
			data.DataType = enums.GetFolderDataType(*dt).Name()
		}
		if *dt != constants.TimeSequenceType {
			data.Fields = base.Filter(unsPo.GetFields(), func(e *types.FieldDefine) bool {
				return !e.IsSystemField()
			})
		} else {
			data.Fields = unsPo.GetFields()
		}
	} else {
		data.Fields = unsPo.GetFields()
	}

	if protocol := unsPo.GetProtocolMap(); len(protocol) > 0 {
		frequency := protocol["frequency"]
		if frequency != nil {
			data.Frequency = fmt.Sprint(frequency)
		}
		if base.P2v(unsPo.GetDataType()) == constants.JsonbType {
			// jsonb 类型不需要导出 fields,但需要导出 jsonFields 作为 fields
			if jsf, has := protocol["jsf"]; has {
				if str, isStr := jsf.(string); isStr {
					JsonUtil.FromJson(str, &data.Fields)
				} else {
					JsonUtil.FromJson(fmt.Sprint(jsf), &data.Fields)
				}
			}
		}
	}

	data.Description = unsPo.GetDescription()
	if pdt := unsPo.GetParentDataType(); pdt != nil {
		data.ParentDataType = enums.GetFolderDataType(*pdt).Name()
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
