package enums

// SysModuleEnum represents system modules
type SysModuleEnum string

const (
	SysModuleAlarm   SysModuleEnum = "system.module.alarm"
	SysModuleUnknown SysModuleEnum = "unknown"
)

// String returns the string representation
func (sm SysModuleEnum) String() string {
	return string(sm)
}

// SysModuleParse returns SysModuleEnum from code
func SysModuleParse(code string) SysModuleEnum {
	switch SysModuleEnum(code) {
	case SysModuleAlarm:
		return SysModuleAlarm
	default:
		return SysModuleUnknown
	}
}
