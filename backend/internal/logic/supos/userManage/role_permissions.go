package userManage

import (
	"sort"
	"strconv"
	"strings"

	"backend/internal/common/enums"
	"backend/internal/repo/relationDB"
	"backend/internal/types"
)

const allButtonsPermission = "button:*"

func buildRolePermissionLists(
	allRows []relationDB.SuposResource,
	assignedRows []relationDB.SuposResource,
) ([]types.RoleResource, []types.RoleResource) {
	pageMap := make(map[string]types.RoleResource)
	allButtons := make(map[string]types.RoleResource)
	selectedButtons := make(map[string]types.RoleResource)

	for _, row := range allRows {
		if row.Type != 3 {
			continue
		}
		code := strings.TrimSpace(row.Code)
		if code == "" {
			continue
		}
		uri := "button:" + code
		allButtons[uri] = types.RoleResource{
			ResourceID: strconv.FormatInt(row.ID, 10),
			URI:        uri,
		}
	}

	for _, row := range assignedRows {
		switch row.Type {
		case 2, 5:
			uri := buildPagePermission(row)
			if uri == "" {
				continue
			}
			pageMap[uri] = types.RoleResource{
				ResourceID: strconv.FormatInt(row.ID, 10),
				URI:        uri,
				Methods:    defaultMethods(),
			}
		case 3:
			code := strings.TrimSpace(row.Code)
			if code == "" {
				continue
			}
			uri := "button:" + code
			selectedButtons[uri] = types.RoleResource{
				ResourceID: strconv.FormatInt(row.ID, 10),
				URI:        uri,
			}
		}
	}

	allow := make([]types.RoleResource, 0, len(pageMap)+1)
	for _, item := range pageMap {
		allow = append(allow, item)
	}

	deny := make([]types.RoleResource, 0)
	if len(selectedButtons) > 0 {
		allow = append(allow, types.RoleResource{URI: allButtonsPermission})
		for uri, item := range allButtons {
			if _, ok := selectedButtons[uri]; ok {
				continue
			}
			deny = append(deny, item)
		}
	}

	sortRoleResources(allow)
	sortRoleResources(deny)
	return allow, deny
}

func mapRoleResourcesToIDsFromRows(
	rows []relationDB.SuposResource,
	allowResources []types.RoleResource,
	denyResources []types.RoleResource,
) []int64 {
	pageIndex := make(map[string]int64)
	buttonIndex := make(map[string]int64)

	for _, row := range rows {
		switch row.Type {
		case 2, 5:
			if uri := buildPagePermission(row); uri != "" {
				pageIndex[uri] = row.ID
			}
		case 3:
			if code := strings.TrimSpace(row.Code); code != "" {
				buttonIndex["button:"+code] = row.ID
			}
		}
	}

	result := make([]int64, 0, len(allowResources))
	seen := make(map[int64]struct{}, len(allowResources))
	deniedButtons := make(map[string]struct{}, len(denyResources))
	allowAllButtons := false

	for _, resource := range denyResources {
		uri := normalizePermissionURI(resource.URI)
		if strings.HasPrefix(uri, "button:") {
			deniedButtons[uri] = struct{}{}
		}
	}

	addID := func(resourceID int64) {
		if resourceID == 0 {
			return
		}
		if _, exists := seen[resourceID]; exists {
			return
		}
		seen[resourceID] = struct{}{}
		result = append(result, resourceID)
	}

	for _, resource := range allowResources {
		uri := normalizePermissionURI(resource.URI)
		if uri == "" || enums.IsDefaultCommonURI(uri) {
			continue
		}
		if uri == allButtonsPermission {
			allowAllButtons = true
			continue
		}
		if strings.HasPrefix(uri, "button:") {
			if _, denied := deniedButtons[uri]; denied {
				continue
			}
			addID(buttonIndex[uri])
			continue
		}
		addID(pageIndex[uri])
	}

	if allowAllButtons {
		buttonURIs := make([]string, 0, len(buttonIndex))
		for uri := range buttonIndex {
			buttonURIs = append(buttonURIs, uri)
		}
		sort.Strings(buttonURIs)
		for _, uri := range buttonURIs {
			if _, denied := deniedButtons[uri]; denied {
				continue
			}
			addID(buttonIndex[uri])
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

func sortRoleResources(items []types.RoleResource) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].URI == items[j].URI {
			return items[i].ResourceID < items[j].ResourceID
		}
		return items[i].URI < items[j].URI
	})
}
