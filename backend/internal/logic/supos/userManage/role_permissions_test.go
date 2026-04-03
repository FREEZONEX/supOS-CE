package userManage

import (
	"testing"

	"backend/internal/repo/relationDB"
	"backend/internal/types"
)

func TestBuildRolePermissionLists(t *testing.T) {
	urlType := 1
	urlAccount := "/account-management"
	allRows := []relationDB.SuposResource{
		{ID: 1, Type: 2, Code: "account-management", URLType: &urlType, URL: &urlAccount},
		{ID: 11, Type: 3, Code: "account.add"},
		{ID: 12, Type: 3, Code: "account.edit"},
		{ID: 13, Type: 3, Code: "account.delete"},
	}
	assignedRows := []relationDB.SuposResource{
		{ID: 1, Type: 2, Code: "account-management", URLType: &urlType, URL: &urlAccount},
		{ID: 11, Type: 3, Code: "account.add"},
		{ID: 12, Type: 3, Code: "account.edit"},
	}

	allow, deny := buildRolePermissionLists(allRows, assignedRows)

	if len(allow) != 2 {
		t.Fatalf("allow length = %d, want 2", len(allow))
	}
	if allow[0].URI != "/account-management" {
		t.Fatalf("allow[0].URI = %q, want /account-management", allow[0].URI)
	}
	if allow[1].URI != allButtonsPermission {
		t.Fatalf("allow[1].URI = %q, want %q", allow[1].URI, allButtonsPermission)
	}
	if len(deny) != 1 || deny[0].URI != "button:account.delete" {
		t.Fatalf("deny = %#v, want only button:account.delete", deny)
	}
}

func TestMapRoleResourcesToIDsFromRowsWithWildcardAndDeny(t *testing.T) {
	urlType := 1
	urlAccount := "/account-management"
	rows := []relationDB.SuposResource{
		{ID: 1, Type: 2, Code: "account-management", URLType: &urlType, URL: &urlAccount},
		{ID: 11, Type: 3, Code: "account.add"},
		{ID: 12, Type: 3, Code: "account.edit"},
		{ID: 13, Type: 3, Code: "account.delete"},
	}

	ids := mapRoleResourcesToIDsFromRows(
		rows,
		[]types.RoleResource{
			{URI: "/account-management"},
			{URI: allButtonsPermission},
		},
		[]types.RoleResource{
			{URI: "button:account.delete"},
		},
	)

	expected := []int64{1, 11, 12}
	if len(ids) != len(expected) {
		t.Fatalf("ids length = %d, want %d, ids = %#v", len(ids), len(expected), ids)
	}
	for i := range expected {
		if ids[i] != expected[i] {
			t.Fatalf("ids[%d] = %d, want %d, ids = %#v", i, ids[i], expected[i], ids)
		}
	}
}

func TestMapRoleResourcesToIDsFromRowsWithExplicitButtons(t *testing.T) {
	rows := []relationDB.SuposResource{
		{ID: 11, Type: 3, Code: "account.add"},
		{ID: 12, Type: 3, Code: "account.edit"},
	}

	ids := mapRoleResourcesToIDsFromRows(
		rows,
		[]types.RoleResource{{URI: "button:account.edit"}},
		nil,
	)

	if len(ids) != 1 || ids[0] != 12 {
		t.Fatalf("ids = %#v, want [12]", ids)
	}
}
