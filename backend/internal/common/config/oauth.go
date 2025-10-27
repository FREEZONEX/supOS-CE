package config

import "net/url"

// OAuthKeyCloakConfig represents Keycloak OAuth configuration values loaded from yaml/env.
type OAuthKeyCloakConfig struct {
	Realm                  string `mapstructure:"realm"`
	ClientName             string `mapstructure:"client_name"`
	ClientID               string `mapstructure:"client_id"`
	ClientSecret           string `mapstructure:"client_secret"`
	AuthorizationGrantType string `mapstructure:"authorization_grant_type"`
	RedirectURI            string `mapstructure:"redirect_uri"`
	IssuerURI              string `mapstructure:"issuer_uri"`
	SuposHome              string `mapstructure:"supos_home"`
	RefreshTokenTime       int64  `mapstructure:"refresh_token_time"`
	SuposClientID          string `mapstructure:"supos_client_id"`
}

// GetRedirectURI returns the redirect URI with default ports stripped.
func (o *OAuthKeyCloakConfig) GetRedirectURI() string {
	return removePortIfDefault(o.RedirectURI)
}

// removePortIfDefault strips default http/https ports from a URL.
func removePortIfDefault(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	port := u.Port()
	isHttpDefault := u.Scheme == "http" && port == "80"
	isHttpsDefault := u.Scheme == "https" && port == "443"

	if isHttpDefault || isHttpsDefault {
		u.Host = u.Hostname()
		return u.String()
	}

	return rawURL
}
