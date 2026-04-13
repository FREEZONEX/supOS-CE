package config

import "testing"

func TestNormalizeLoginPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty uses default", in: "", want: DefaultLoginPath},
		{name: "without leading slash", in: "tier0-login", want: "/tier0-login"},
		{name: "single leading slash", in: "/tier0-login", want: "/tier0-login"},
		{name: "double leading slash", in: "//tier0-login", want: "/tier0-login"},
		{name: "keeps query", in: "//tier0-login?source=oauth", want: "/tier0-login?source=oauth"},
		{name: "keeps fragment", in: "tier0-login#form", want: "/tier0-login#form"},
		{name: "absolute url unchanged", in: "https://example.com/login", want: "https://example.com/login"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeLoginPath(tt.in); got != tt.want {
				t.Fatalf("NormalizeLoginPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
