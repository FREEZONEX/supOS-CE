package vo

import (
	"backend/internal/types"
	"fmt"
	"time"

	"backend/internal/common/constants"
	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
)

// FieldDefineVo represents a field definition view object.
type FieldDefineVo struct {
	Name        string `json:"name" validate:"required"` // 字段名
	Type        string `json:"type" validate:"required"` // 字段类型：int, long, float, string, boolean
	Unique      bool   `json:"unique,omitzero"`          // 是否唯一约束
	Index       string `json:"index,omitzero"`           // modeBus 协议时字段对应的数组下标
	System      bool   `json:"system,omitzero"`          // 是否系统预置字段
	DisplayName string `json:"displayName,omitzero"`     // 显式名
	Remark      string `json:"remark,omitzero"`          // 备注
	MaxLen      string `json:"maxLen,omitzero"`          // 最大长度
}

// NewFieldDefineVo creates a FieldDefineVo from FieldDefine.
func NewFieldDefineVo(bo *types.FieldDefine) *FieldDefineVo {
	vo := &FieldDefineVo{
		Name: bo.Name,
		Type: bo.Type,
	}
	if idx := bo.Index; idx != nil {
		vo.Index = *idx
	}
	if bo.Unique != nil {
		vo.Unique = *bo.Unique
	}
	return vo
}

// NewFieldDefineVoSimple creates a FieldDefineVo with basic fields.
func NewFieldDefineVoSimple(name, fieldType string) *FieldDefineVo {
	return &FieldDefineVo{
		Name: name,
		Type: fieldType,
	}
}

// NewFieldDefineVoWithUnique creates a FieldDefineVo with unique constraint.
func NewFieldDefineVoWithUnique(name, fieldType string, unique bool) *FieldDefineVo {
	return &FieldDefineVo{
		Name:   name,
		Type:   fieldType,
		Unique: unique,
	}
}

// NewFieldDefineVoWithIndex creates a FieldDefineVo with index.
func NewFieldDefineVoWithIndex(name, fieldType, index string) *FieldDefineVo {
	return &FieldDefineVo{
		Name:  name,
		Type:  fieldType,
		Index: index,
	}
}

// GetSystem returns the system field status, auto-detecting from name prefix if not set.
func (v *FieldDefineVo) GetSystem() bool {
	if !v.System && v.Name != "" && len(v.Name) >= len(constants.SystemFieldPrev) {
		if v.Name[:len(constants.SystemFieldPrev)] == constants.SystemFieldPrev {
			return true
		}
	}
	return v.System
}

// IsUnique checks if the field has a unique constraint.
func (v *FieldDefineVo) IsUnique() bool {
	return v.Unique
}

// Convert converts the VO to FieldDefine DTO.
func (v *FieldDefineVo) Convert() *types.FieldDefine {
	return vo2bo(v)
}

// ConvertArray converts an array of FieldDefineVo to FieldDefine array.
func ConvertArray(vfs []*FieldDefineVo) []*types.FieldDefine {
	if len(vfs) == 0 {
		return nil
	}
	fs := make([]*types.FieldDefine, len(vfs))
	for i, vo := range vfs {
		fs[i] = vo2bo(vo)
	}
	return fs
}

func vo2bo(vo *FieldDefineVo) *types.FieldDefine {
	fieldType, _ := types.GetFieldTypeByName(vo.Type)
	define := &types.FieldDefine{
		Name:  vo.Name,
		Type:  fieldType.Name(),
		Index: &vo.Index,
	}
	define.Unique = &vo.Unique
	define.DisplayName = &vo.DisplayName
	define.Remark = &vo.Remark
	return define
}

// HistoryValueVo represents a history query result response object.
type HistoryValueVo struct {
	Results            []*HistoryResult `json:"results,omitzero"`             // 查询结果集合，每个元素代表一个文件的聚合查询结果
	Unauthorized       []string         `json:"unauthorized,omitzero"`        // 无授权的文件别名集合
	NotExistAttributes []string         `json:"notExsistAtrributes,omitzero"` // 不存在的文件别名集合
}

// HistoryResult represents a single file's aggregated query result.
type HistoryResult struct {
	Alias    string     `json:"alias"`    // 文件别名
	Function string     `json:"function"` // 聚合函数，如 first、sum、mean 等
	HasNext  bool       `json:"hasNext"`  // 是否有下一页数据
	Fields   []string   `json:"fields"`   // 字段顺序，返回的数据字段顺序与此保持一致
	Datas    [][]string `json:"datas"`    // 数据集合，每条记录为一个 List，字段顺序与 fields 对应
}

// LabelVo represents a label view object.
type LabelVo struct {
	ID        string    `json:"id,omitzero"`        // 标签ID：已有标签时必传，新建标签时不需要传
	LabelName string    `json:"labelName,omitzero"` // 标签名称，新建标签时，必传
	CreateAt  time.Time `json:"createAt,omitzero"`  // 创建时间
}

// UserAttributeVo represents user attributes view object.
type UserAttributeVo struct {
	FirstTimeLogin int    `json:"firstTimeLogin"`  // 是否首次登录 (1: 是, 0: 否)
	TipsEnable     int    `json:"tipsEnable"`      // 是否开启tips (1: 是, 0: 否)
	HomePage       string `json:"homePage"`        // 用户自定义首页
	Phone          string `json:"phone,omitzero"`  // 手机号
	Source         string `json:"source,omitzero"` // 用户来源
}

// NewUserAttributeVo creates a UserAttributeVo with default values.
func NewUserAttributeVo() *UserAttributeVo {
	return &UserAttributeVo{
		FirstTimeLogin: 1,
		TipsEnable:     1,
		HomePage:       constants.DefaultHomepage,
	}
}

// UserInfoVo represents user information view object.
type UserInfoVo struct {
	UserAttributeVo                          // Embedded attributes
	Sub               string                 `json:"sub"`                       // 用户的唯一标识符（用户ID）
	PreferredUsername string                 `json:"preferredUsername"`         // 用户名
	Email             string                 `json:"email,omitzero"`            // 邮箱
	EmailVerified     bool                   `json:"email_verified,omitzero"`   // 邮箱验证状态
	FirstName         string                 `json:"firstName,omitzero"`        // 名字
	Enabled           bool                   `json:"enabled,omitzero"`          // 是否启用
	RoleList          []*authdto.RoleDto     `json:"roleList,omitzero"`         // 角色列表
	ResourceList      []*authdto.ResourceDto `json:"resourceList,omitzero"`     // 资源列表
	DenyResourceList  []*authdto.ResourceDto `json:"denyResourceList,omitzero"` // 拒绝策略资源列表
	MainLanguage      string                 `json:"mainLanguage,omitzero"`     // 主语言
	SuperAdmin        bool                   `json:"superAdmin,optional"`       // 是否为超级管理员
}

// NewUserInfoVo creates a UserInfoVo with basic info.
func NewUserInfoVo(sub, preferredUsername string) *UserInfoVo {
	return &UserInfoVo{
		UserAttributeVo:   *NewUserAttributeVo(),
		Sub:               sub,
		PreferredUsername: preferredUsername,
	}
}

// IsSuperAdmin checks if the user is a super admin.
func (u *UserInfoVo) IsSuperAdmin() bool {
	if len(u.RoleList) == 0 {
		return false
	}
	for _, role := range u.RoleList {
		if role.RoleID == enums.RoleSuperAdmin.ID {
			return true
		}
	}
	return false
}

// Guest creates a guest user.
func Guest() *UserInfoVo {
	user := NewUserInfoVo("guest", "guest")
	user.FirstName = "guest"
	user.Enabled = true
	user.FirstTimeLogin = 0
	user.TipsEnable = 0
	user.HomePage = constants.DefaultHomepage
	return user
}

// String returns a simple JSON-like string representation.
func (u *UserInfoVo) String() string {
	// A more complete implementation might use json.Marshal, but this matches the Java version's intent for simple logging.
	return fmt.Sprintf(`{"sub":"%s", "preferredUsername":"%s"}`, u.Sub, u.PreferredUsername)
}

// UserManageVo represents user management view object.
type UserManageVo struct {
	UserAttributeVo                      // Embedded attributes
	ID                string             `json:"id"`                         // 用户ID
	Email             string             `json:"email,omitzero"`             // 邮箱
	EmailVerified     bool               `json:"emailVerified,omitzero"`     // 邮箱验证状态
	FirstName         string             `json:"firstName,omitzero"`         // 名字
	PreferredUsername string             `json:"preferredUsername,omitzero"` // 用户的首选用户名
	Sub               string             `json:"sub,omitzero"`               // 用户的唯一标识符：用户ID
	Enabled           bool               `json:"enabled,omitzero"`           // 是否启用
	RoleList          []*authdto.RoleDto `json:"roleList,omitzero"`          // 角色列表
}

// NewUserManageVo creates a UserManageVo with basic info.
func NewUserManageVo(id, preferredUsername string) *UserManageVo {
	return &UserManageVo{
		UserAttributeVo:   *NewUserAttributeVo(),
		ID:                id,
		PreferredUsername: preferredUsername,
	}
}

// String returns a simple JSON-like string representation.
func (u *UserManageVo) String() string {
	return fmt.Sprintf(`{"id":"%s", "preferredUsername":"%s"}`, u.ID, u.PreferredUsername)
}
