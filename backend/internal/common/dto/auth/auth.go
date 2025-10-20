package auth

type AccessTokenDto struct {
	AccessToken      string `json:"accessToken"`
	ExpiresIn        int    `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
}

// AddUserDto uses validator struct tags to enforce validation rules.
type AddUserDto struct {
	ID        string    `json:"id,omitzero"`
	Username  string    `json:"username" validate:"required,min=3"`
	Password  string    `json:"password" validate:"required,min=3,max=10"` // Simplified regex
	Enabled   bool      `json:"enabled"`
	Email     string    `json:"email" validate:"omitzero,email"`
	FirstName string    `json:"firstName,omitzero"`
	Phone     string    `json:"phone,omitzero"`
	Source    string    `json:"source,omitzero"`
	RoleList  []RoleDto `json:"roleList,omitzero"`
}

type RoleDto struct {
	RoleID          string `json:"roleId"`
	RoleName        string `json:"roleName"`
	RoleDescription string `json:"roleDescription,omitzero"`
	ClientRole      bool   `json:"clientRole,omitzero"`
}

type ResourceDto struct {
	PolicyID   string   `json:"policyId"`
	ResourceID string   `json:"resourceId"`
	URI        string   `json:"uri"`
	Methods    []string `json:"methods"`
}

type UpdateUserDto struct {
	UserID         string    `json:"userId" validate:"required"`
	Username       string    `json:"username,omitzero"`
	Password       string    `json:"password,omitzero" validate:"omitzero,min=3,max=10"`
	Enabled        bool      `json:"enabled,omitzero"`
	Email          string    `json:"email,omitzero" validate:"omitzero,email"`
	FirstName      string    `json:"firstName,omitzero"`
	Phone          string    `json:"phone,omitzero"`
	FirstTimeLogin int       `json:"firstTimeLogin,omitzero"`
	TipsEnable     int       `json:"tipsEnable,omitzero"`
	HomePage       string    `json:"homePage,omitzero"`
	RoleList       []RoleDto `json:"roleList,omitzero"`
	OperateRole    bool      `json:"operateRole,omitzero"`
	Source         string    `json:"source,omitzero"`
}

type UpdateRoleDto struct {
	UserID   string    `json:"userId" validate:"required"`
	Type     int       `json:"type,omitzero"` // 1: set, 2: unset
	RoleList []RoleDto `json:"roleList"`
}

// RoleSaveDto
// TODO: Group validation (Create vs Update) and custom @RoleNameValidator
// need to be implemented in the service layer logic.
type RoleSaveDto struct {
	ID                string        `json:"id"` // Required on update
	Name              string        `json:"name" validate:"required"`
	DenyResourceList  []ResourceDto `json:"denyResourceList"`
	AllowResourceList []ResourceDto `json:"allowResourceList"`
}

type ResetPasswordDto struct {
	UserID      string `json:"userId" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type UserQueryDto struct {
	PageNo            int64  `json:"pageNo" form:"pageNo"`
	PageSize          int64  `json:"pageSize" form:"pageSize"`
	PreferredUsername string `json:"preferredUsername,omitzero"`
	FirstName         string `json:"firstName,omitzero"`
	Email             string `json:"email,omitzero"`
	RoleName          string `json:"roleName,omitzero"`
	Enabled           bool   `json:"enabled"`
}
