package auth

import "backend/internal/common/dto"

type KeycloakUserInfoDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type KeycloakRoleInfoDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitzero"`
}

type KeycloakResourceInfoDto struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	TypeName    string `json:"typeName,omitzero"`
	DisplayName string `json:"displayName,omitzero"`
}

type KeycloakPolicyInfoDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitzero"`
}

type KeycloakCreateUserDto struct {
	ID         string                `json:"id,omitzero"`
	Username   string                `json:"username"`
	Enabled    bool                  `json:"enabled"`
	Email      string                `json:"email,omitzero"`
	FirstName  string                `json:"firstName,omitzero"`
	Attributes *dto.UserAttributeDto `json:"attributes,omitzero"`
}
