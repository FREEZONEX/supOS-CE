package enums

// GlobalExportModuleEnum represents global export module types
type GlobalExportModuleEnum struct {
	Code string
	Name string
}

var (
	GlobalExportModuleUNS = GlobalExportModuleEnum{
		Code: "uns",
		Name: "",
	}

	GlobalExportModuleSourceFlow = GlobalExportModuleEnum{
		Code: "sourceFlow",
		Name: "",
	}

	GlobalExportModuleEventFlow = GlobalExportModuleEnum{
		Code: "eventFlow",
		Name: "",
	}

	GlobalExportModuleDashboard = GlobalExportModuleEnum{
		Code: "dashboard",
		Name: "",
	}
)

// All global export modules
var allGlobalExportModules = []GlobalExportModuleEnum{
	GlobalExportModuleUNS,
	GlobalExportModuleSourceFlow,
	GlobalExportModuleEventFlow,
	GlobalExportModuleDashboard,
}

// Is checks if the code matches
func (gem GlobalExportModuleEnum) Is(code string) bool {
	return gem.Code == code
}

// GlobalExportModuleFromCode returns GlobalExportModuleEnum from code
func GetGlobalExportModuleFromCode(code string) *GlobalExportModuleEnum {
	for _, module := range allGlobalExportModules {
		if module.Code == code {
			return &module
		}
	}
	return nil
}
