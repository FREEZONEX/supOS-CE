package common

import "testing"

func TestFormatWorkspaceAPIKey(t *testing.T) {
	key := FormatWorkspaceAPIKey("sk-svc-", 10001, "abc_DEF")
	if key != "sk-svc-ws7pt_abc_DEF" {
		t.Fatalf("FormatWorkspaceAPIKey() = %q", key)
	}
	if got := ParseWorkspaceIDFromAPIKey(key); got != 10001 {
		t.Fatalf("ParseWorkspaceIDFromAPIKey() = %d, want 10001", got)
	}
}

func TestParseWorkspaceIDFromAPIKeyInvalid(t *testing.T) {
	for _, key := range []string{
		"",
		"t0_legacy",
		"sk-svc-token",
		"sk-svc-ws_token",
		"sk-svc-ws0_token",
	} {
		if got := ParseWorkspaceIDFromAPIKey(key); got != 0 {
			t.Fatalf("ParseWorkspaceIDFromAPIKey(%q) = %d, want 0", key, got)
		}
	}
}
