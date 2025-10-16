package enums

// ParentResourceEnum represents parent resource navigation items
type ParentResourceEnum struct {
	ID        int64
	Code      string
	Comment   string
	GroupType int // 菜单分组 1-导航 2-菜单 3-tab
}

var (
	ParentResourceNavAppspace = ParentResourceEnum{
		ID:        4,
		Code:      "menu.tag.appspace",
		Comment:   "应用集",
		GroupType: 1,
	}

	ParentResourceNavDevtools = ParentResourceEnum{
		ID:        50,
		Code:      "menu.tag.devtools",
		Comment:   "工具集",
		GroupType: 1,
	}

	ParentResourceNavSystem = ParentResourceEnum{
		ID:        60,
		Code:      "menu.tag.system",
		Comment:   "系统管理",
		GroupType: 1,
	}

	ParentResourceNavUNS = ParentResourceEnum{
		ID:        60,
		Code:      "menu.tag.uns",
		Comment:   "数据管理",
		GroupType: 1,
	}
)

// All parent resources
var allParentResources = map[string]ParentResourceEnum{
	ParentResourceNavAppspace.Code: ParentResourceNavAppspace,
	ParentResourceNavDevtools.Code: ParentResourceNavDevtools,
	ParentResourceNavSystem.Code:   ParentResourceNavSystem,
	ParentResourceNavUNS.Code:      ParentResourceNavUNS,
}

// GetParentResourceFromCode returns ParentResourceEnum from code
func GetParentResourceFromCode(code string) (ParentResourceEnum, bool) {
	res, ok := allParentResources[code]
	return res, ok
}
