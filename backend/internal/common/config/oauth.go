package config

import "net/url"

// OAuthKeyCloakConfig represents Keycloak OAuth configuration
type OAuthKeyCloakConfig struct {
	Realm                  string `mapstructure:"realm"`                    // Keycloak realm
	ClientName             string `mapstructure:"client_name"`              // 客户端名称
	ClientID               string `mapstructure:"client_id"`                // 客户端ID
	ClientSecret           string `mapstructure:"client_secret"`            // 客户端密钥
	AuthorizationGrantType string `mapstructure:"authorization_grant_type"` // 授权类型
	RedirectURI            string `mapstructure:"redirect_uri"`             // 重定向URI
	IssuerURI              string `mapstructure:"issuer_uri"`               // 发行者URI
	SuposHome              string `mapstructure:"supos_home"`               // supOS主页
	RefreshTokenTime       int64  `mapstructure:"refresh_token_time"`       // Token刷新时间
}

// GetRedirectURI returns the redirect URI with default port removed
func (o *OAuthKeyCloakConfig) GetRedirectURI() string {
	return removePortIfDefault(o.RedirectURI)
}

// removePortIfDefault uses the standard library's URL parser to remove default ports.
// This approach is more robust and safer than using regular expressions.
func removePortIfDefault(rawURL string) string {
	// Attempt to parse the URL. If it's malformed, we can't safely modify it,
	// so we return the original string.
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Return original on parsing error
	}

	// The Port() method correctly extracts the port number as a string.
	port := u.Port()

	// Check if the scheme and port match the default combinations.
	isHttpDefault := u.Scheme == "http" && port == "80"
	isHttpsDefault := u.Scheme == "https" && port == "443"

	if isHttpDefault || isHttpsDefault {
		// The Host field of url.URL contains both hostname and port (e.g., "example.com:80").
		// To remove the port, we can simply set the Host to be just the hostname.
		u.Host = u.Hostname()
		// Reconstruct the URL string from the modified parts.
		return u.String()
	}

	// If no modification was needed, return the original URL.
	return rawURL
}
