package iam

import (
	"sort"
	"strconv"
	"strings"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/repo/relationDB"

	"gorm.io/gorm"
)

var supportedMethods = []string{"get", "post", "put", "delete", "patch", "head", "options"}

func defaultMethods() []string {
	result := make([]string, len(supportedMethods))
	copy(result, supportedMethods)
	return result
}

func expandPermissionRows(rows []relationDB.SuposResource) ([]*authdto.ResourceDto, []string, error) {
	pageMap := make(map[string]*authdto.ResourceDto, len(rows))
	buttonSet := make(map[string]struct{})
	for _, row := range rows {
		switch row.Type {
		case 2, 5:
			uri := buildPagePermission(row)
			if strings.TrimSpace(uri) == "" {
				continue
			}
			pageMap[uri] = &authdto.ResourceDto{
				ResourceID: stringifyInt64(row.ID),
				URI:        uri,
				Methods:    defaultMethods(),
			}
		case 3:
			code := strings.TrimSpace(row.Code)
			if code == "" {
				continue
			}
			buttonSet["button:"+code] = struct{}{}
		}
	}

	pages := make([]*authdto.ResourceDto, 0, len(pageMap))
	for _, page := range pageMap {
		pages = append(pages, page)
	}
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].URI < pages[j].URI
	})

	buttons := make([]string, 0, len(buttonSet))
	for button := range buttonSet {
		buttons = append(buttons, button)
	}
	sort.Strings(buttons)
	return pages, buttons, nil
}

func mapResourcesToIDs(db *gorm.DB, resources []*authdto.ResourceDto) ([]int64, error) {
	if len(resources) == 0 {
		return nil, nil
	}
	var rows []relationDB.SuposResource
	if err := db.Model(&relationDB.SuposResource{}).
		Where("type IN ?", []int{2, 3, 5}).
		Where("COALESCE(enable, true) = true").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	uriToID := make(map[string]int64, len(rows)*2)
	for _, row := range rows {
		switch row.Type {
		case 2, 5:
			if uri := buildPagePermission(row); uri != "" {
				uriToID[uri] = row.ID
			}
		case 3:
			if code := strings.TrimSpace(row.Code); code != "" {
				uriToID["button:"+code] = row.ID
			}
		}
	}

	result := make([]int64, 0, len(resources))
	seen := make(map[int64]struct{}, len(resources))
	for _, resource := range resources {
		if resource == nil {
			continue
		}
		uri := normalizePermissionURI(resource.URI)
		if uri == "" || enums.IsDefaultCommonURI(uri) {
			continue
		}
		resourceID, ok := uriToID[uri]
		if !ok {
			continue
		}
		if _, exists := seen[resourceID]; exists {
			continue
		}
		seen[resourceID] = struct{}{}
		result = append(result, resourceID)
	}
	return result, nil
}

func normalizeRoleInputs(inputs []RoleSyncInput) []RoleSyncInput {
	if len(inputs) == 0 {
		return []RoleSyncInput{
			{
				Role: &authdto.RoleDto{
					RoleID:          enums.RoleNormalUser.ID,
					RoleName:        "user",
					RoleDescription: "builtin user role",
					ClientRole:      true,
				},
			},
		}
	}

	result := make([]RoleSyncInput, 0, len(inputs))
	for _, input := range inputs {
		role := normalizeRole(input.Role)
		if role == nil {
			continue
		}
		result = append(result, RoleSyncInput{
			Role:      role,
			Resources: input.Resources,
		})
	}
	if len(result) == 0 {
		return normalizeRoleInputs(nil)
	}
	return result
}

func normalizeRole(role *authdto.RoleDto) *authdto.RoleDto {
	if role == nil {
		return nil
	}
	roleName := strings.TrimSpace(role.RoleName)
	roleID := strings.TrimSpace(role.RoleID)
	description := strings.TrimSpace(role.RoleDescription)
	switch strings.ToLower(roleName) {
	case "super-admin", "admin":
		return &authdto.RoleDto{
			RoleID:          enums.RoleSuperAdmin.ID,
			RoleName:        "admin",
			RoleDescription: firstNonEmpty(description, "builtin admin role"),
			ClientRole:      true,
		}
	case "normal-user", "user":
		return &authdto.RoleDto{
			RoleID:          enums.RoleNormalUser.ID,
			RoleName:        "user",
			RoleDescription: firstNonEmpty(description, "builtin user role"),
			ClientRole:      true,
		}
	default:
		if roleID == "" {
			return nil
		}
		return &authdto.RoleDto{
			RoleID:          roleID,
			RoleName:        firstNonEmpty(roleName, roleID),
			RoleDescription: description,
			ClientRole:      true,
		}
	}
}

func normalizeRoleKey(role *authdto.RoleDto) string {
	if role == nil {
		return ""
	}
	if strings.EqualFold(role.RoleID, enums.RoleSuperAdmin.ID) || strings.EqualFold(role.RoleName, "admin") {
		return "admin"
	}
	if strings.EqualFold(role.RoleID, enums.RoleNormalUser.ID) || strings.EqualFold(role.RoleName, "user") {
		return "user"
	}
	return strings.ToLower(strings.TrimSpace(role.RoleName))
}

func isBuiltinRole(role *authdto.RoleDto) bool {
	key := normalizeRoleKey(role)
	return key == "admin" || key == "user"
}

func buildPagePermission(resource relationDB.SuposResource) string {
	code := strings.TrimSpace(resource.Code)
	if resource.URLType != nil && *resource.URLType == 1 && resource.URL != nil && strings.TrimSpace(*resource.URL) != "" {
		return strings.TrimSpace(*resource.URL)
	}
	if code == "" {
		return ""
	}
	return "/" + code
}

func normalizePermissionURI(uri string) string {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return ""
	}
	if idx := strings.Index(uri, "$"); idx >= 0 {
		return strings.TrimSpace(uri[:idx])
	}
	return uri
}

func stringifyInt64(val int64) string {
	return strconv.FormatInt(val, 10)
}
