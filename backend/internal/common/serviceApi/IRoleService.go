package serviceApi

import authdto "backend/internal/common/dto/auth"

type IRoleService interface {
	GetRoleListByUserId(userID string) ([]*authdto.RoleDto, error)
}
