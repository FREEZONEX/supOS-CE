package open

import (
	"testing"
	"time"

	"backend/internal/types"
)

func TestParseGetMembersUpdatedAtRange(t *testing.T) {
	req := &types.GetMembersReq{
		UpdatedAtStart: "2026-04-01T00:00:00Z",
		UpdatedAtEnd:   "2026-04-30T23:59:59Z",
	}

	start, end, err := parseGetMembersUpdatedAtRange(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if start == nil || !start.Equal(time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected start: %v", start)
	}
	if end == nil || !end.Equal(time.Date(2026, 4, 30, 23, 59, 59, 0, time.UTC)) {
		t.Fatalf("unexpected end: %v", end)
	}
}

func TestParseGetMembersUpdatedAtRangeInvalidStart(t *testing.T) {
	_, _, err := parseGetMembersUpdatedAtRange(&types.GetMembersReq{UpdatedAtStart: "bad-time"})
	if err == nil {
		t.Fatal("expected invalid start error")
	}
}

func TestParseGetMembersUpdatedAtRangeInvalidEnd(t *testing.T) {
	_, _, err := parseGetMembersUpdatedAtRange(&types.GetMembersReq{UpdatedAtEnd: "bad-time"})
	if err == nil {
		t.Fatal("expected invalid end error")
	}
}
