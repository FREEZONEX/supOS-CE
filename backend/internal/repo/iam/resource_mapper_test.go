package iam

import (
	"backend/internal/repo/relationDB"
	"testing"
)

func TestExpandPermissionRows(t *testing.T) {
	urlType := 1
	url := "/uns"
	rows := []relationDB.SuposResource{
		{ID: 1, Type: 2, Code: "ignored", URLType: &urlType, URL: &url},
		{ID: 2, Type: 5, Code: "dashboards"},
		{ID: 3, Type: 3, Code: "save"},
	}

	pages, buttons, err := expandPermissionRows(rows)
	if err != nil {
		t.Fatalf("expandPermissionRows() error = %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("expandPermissionRows() pages len = %d, want 2", len(pages))
	}
	if pages[0].URI != "/dashboards" || pages[1].URI != "/uns" {
		t.Fatalf("expandPermissionRows() pages = %#v", pages)
	}
	if len(buttons) != 1 || buttons[0] != "button:save" {
		t.Fatalf("expandPermissionRows() buttons = %#v", buttons)
	}
}

func TestNormalizeRoleInputsDefault(t *testing.T) {
	roles := normalizeRoleInputs(nil)
	if len(roles) != 1 || roles[0].Role == nil {
		t.Fatalf("normalizeRoleInputs(nil) = %#v", roles)
	}
	if roles[0].Role.RoleName != "user" {
		t.Fatalf("normalizeRoleInputs(nil) role name = %q, want user", roles[0].Role.RoleName)
	}
}
