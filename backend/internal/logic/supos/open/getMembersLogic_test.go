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

	updatedAtRange, err := parseGetMembersUpdatedAtRange(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedAtRange.Start == nil || !updatedAtRange.Start.Equal(time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected start: %v", updatedAtRange.Start)
	}
	if updatedAtRange.End == nil || !updatedAtRange.End.Equal(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected end: %v", updatedAtRange.End)
	}
	if !updatedAtRange.EndExclusive {
		t.Fatal("expected second-precision end to be exclusive upper bound")
	}
}

func TestParseGetMembersUpdatedAtRangeInvalidStart(t *testing.T) {
	_, err := parseGetMembersUpdatedAtRange(&types.GetMembersReq{UpdatedAtStart: "bad-time"})
	if err == nil {
		t.Fatal("expected invalid start error")
	}
}

func TestParseGetMembersUpdatedAtRangeInvalidEnd(t *testing.T) {
	_, err := parseGetMembersUpdatedAtRange(&types.GetMembersReq{UpdatedAtEnd: "bad-time"})
	if err == nil {
		t.Fatal("expected invalid end error")
	}
}

func TestParseGetMembersUpdatedAtRangeKeepsFractionalEndInclusive(t *testing.T) {
	updatedAtRange, err := parseGetMembersUpdatedAtRange(&types.GetMembersReq{
		UpdatedAtEnd: "2026-05-20T06:11:41.123Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedAtRange.End == nil || !updatedAtRange.End.Equal(time.Date(2026, 5, 20, 6, 11, 41, 123000000, time.UTC)) {
		t.Fatalf("unexpected end: %v", updatedAtRange.End)
	}
	if updatedAtRange.EndExclusive {
		t.Fatal("expected fractional end to stay inclusive")
	}
}

func TestRFC3339HasFractionalSecond(t *testing.T) {
	testCases := []struct {
		value string
		want  bool
	}{
		{value: "2026-05-20T06:11:41Z", want: false},
		{value: "2026-05-20T06:11:41+08:00", want: false},
		{value: "2026-05-20T06:11:41.123Z", want: true},
		{value: "2026-05-20T06:11:41.123+08:00", want: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.value, func(t *testing.T) {
			if got := rfc3339HasFractionalSecond(testCase.value); got != testCase.want {
				t.Fatalf("expected %v, got %v", testCase.want, got)
			}
		})
	}
}
