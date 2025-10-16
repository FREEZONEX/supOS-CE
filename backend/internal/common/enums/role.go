package enums

// RoleEnum represents system roles
type RoleEnum struct {
	ID       string
	Name     string
	I18nCode string
	Comment  string
}

var (
	RoleSuperAdmin = RoleEnum{
		ID:       "7ca9f922-0d35-44cf-8747-8dcfd5e66f8e",
		Name:     "super-admin",
		I18nCode: "user.role.supAdmin",
		Comment:  "超级管理员",
	}

	RoleAdmin = RoleEnum{
		ID:       "a22ce15f-7bef-4e2e-9909-78f51b91c799",
		Name:     "admin",
		I18nCode: "user.role.admin",
		Comment:  "管理员",
	}

	RoleNormalUser = RoleEnum{
		ID:       "71dd6dc2-6b12-4273-9ec0-b44b86e5b500",
		Name:     "normal-user",
		I18nCode: "user.role.normalUser",
		Comment:  "一般用户",
	}
)

// All roles
var AllRoles = []RoleEnum{
	RoleSuperAdmin,
	RoleAdmin,
	RoleNormalUser,
}

// Ignore role IDs (roles to ignore in role list display)
var IgnoreRoleIDs = map[string]bool{
	"625d093d-1333-47d4-92fa-dded93a4f90a": true, // shimu
	"831f62ab-d306-4b11-882e-b23c37ee8c7e": true, // uma_protection
	"2152d19d-e4f9-488d-8509-e49cf239596a": true, // supos-default
	"a22ce15f-7bef-4e2e-9909-78f51b91c799": true, // admin
}

// Ignore role names
var IgnoreRoleNames = map[string]bool{
	"supos-default":    true,
	"ldap-initialized": true,
}

// RoleParse returns RoleEnum from ID
func RoleParse(id string) (RoleEnum, bool) {
	for _, role := range AllRoles {
		if role.ID == id {
			return role, true
		}
	}
	return RoleEnum{}, false
}

// RoleParseName returns RoleEnum from name
func RoleParseName(name string) (RoleEnum, bool) {
	for _, role := range AllRoles {
		if role.Name == name {
			return role, true
		}
	}
	return RoleEnum{}, false
}

// IsIgnoredRoleID checks if role ID should be ignored
func IsIgnoredRoleID(id string) bool {
	return IgnoreRoleIDs[id]
}

// IsIgnoredRoleName checks if role name should be ignored
func IsIgnoredRoleName(name string) bool {
	return IgnoreRoleNames[name]
}
