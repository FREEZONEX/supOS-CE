package clients

import (
	berrors "backend/internal/common/errors"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/zeromicro/go-zero/core/logx"
)

// --- DTO Structs ---

type AccessTokenDto struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type KeycloakUserInfoDto struct {
	ID                string         `json:"id"`
	Username          string         `json:"username"`
	FirstName         string         `json:"firstName"`
	LastName          string         `json:"lastName"`
	Email             string         `json:"email"`
	Enabled           bool           `json:"enabled"`
	Attributes        map[string]any `json:"attributes"`
	PreferredUsername string         `json:"preferredUsername"`
}

type KeycloakRoleInfoDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Composite   bool   `json:"composite"`
	ClientRole  bool   `json:"clientRole"`
	ContainerID string `json:"containerId"`
}

type KeycloakResourceInfoDto struct {
	ID          string   `json:"_id,omitempty"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName,omitempty"`
	Type        string   `json:"type,omitempty"`
	URIs        []string `json:"uris,omitempty"`
}

type KeycloakPolicyInfoDto struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	Type             string `json:"type"`
	Logic            string `json:"logic,omitempty"`
	DecisionStrategy string `json:"decisionStrategy,omitempty"`
}

type keycloakClientRepresentation struct {
	ID       string `json:"id"`
	ClientID string `json:"clientId"`
	Name     string `json:"name,omitempty"`
}

// --- DTO & Error Structs ---

type APIError struct {
	Response *http.Response
	Body     []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed with status %s: %s", e.Response.Status, string(e.Body))
}

// --- Keycloak Client ---

const (
	adminTokenCacheKey = "admin_token"
)

// OAuthKeyCloakConfig represents Keycloak OAuth configuration values loaded from yaml/env.
type KeycloakConfig struct {
	Realm string `json:",optional,env=OAUTH_REALM,default=supos"`
	// YAML: OAuthKeyCloak.client-name
	ClientName string `json:",optional,env=OAUTH_CLIENT_NAME,default=supos"`
	// YAML: OAuthKeyCloak.client-id
	ClientID string `json:",optional,env=OAUTH_CLIENT_ID,default=supos"`
	// YAML: OAuthKeyCloak.client-secret
	ClientSecret string `json:",optional,env=OAUTH_CLIENT_SECRET,default=VaOS2makbDhJJsLlYPt4Wl87bo9VzXiO"`
	// YAML: OAuthKeyCloak.authorization-grant-type
	AuthorizationGrantType string `json:",optional,env=OAUTH_GRANT_TYPE,default=authorization_code"`
	// YAML: OAuthKeyCloak.redirect-uri
	RedirectURI string `json:",optional,env=OAUTH_REDIRECT_URI"`
	// YAML: OAuthKeyCloak.issuer-uri
	IssuerURI string `json:",optional,env=OAUTH_ISSUER_URI"`
	// YAML: OAuthKeyCloak.supos-home
	SuposHome string `json:",optional,env=OAUTH_SUPOS_HOME"`
	// YAML: OAuthKeyCloak.refresh-token-time
	RefreshTokenTime string `json:",optional,env=OAUTH_REFRESH_TOKEN_TIME"`
	// Optional: if provided in YAML as OAuthKeyCloak.supos-client-id; otherwise resolved via admin API.
	SuposClientID string `json:",optional"`
}

// GetRedirectURI returns the redirect URI with default ports stripped.
func (o *KeycloakConfig) GetRedirectURI() string {
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

var Tier0ClientID = "a7b53e5e-3567-470a-9da1-94cc0c7f18e6"

type KeycloakClient struct {
	config     KeycloakConfig
	httpClient *http.Client
	tokenCache *cache.Cache
}

var (
	keycloakClient *KeycloakClient
	keycloakOnce   sync.Once
)

func InitKeycloakClient(config KeycloakConfig) *KeycloakClient {
	fmt.Println("KeycloakConfig", config)
	config.RedirectURI = config.GetRedirectURI()
	if config.ClientID == "" {
		return &KeycloakClient{} //允许 mock
	}
	keycloakOnce.Do(func() {
		kc := &KeycloakClient{
			config:     config,
			httpClient: &http.Client{Timeout: 10 * time.Second},
			tokenCache: cache.New(10*time.Minute, 15*time.Minute),
		}

		// if kc.config.SuposClientID == "" {
		// 	if err := kc.resolveSuposClientID(); err != nil {
		// 		logx.Errorf("failed to resolve Keycloak client uuid for clientId %q: %v", kc.config.ClientID, err)
		// 		// panic(fmt.Errorf("resolve Keycloak client uuid: %w", err))
		// 		fmt.Println("KeycloakConfig init failed!!!")
		// 		keycloakClient = kc
		// 		return
		// 	}
		// }
		kc.config.SuposClientID = Tier0ClientID
		fmt.Println("KeycloakConfig init success!!!")
		keycloakClient = kc
	})
	return keycloakClient
}

func (kc *KeycloakClient) resolveSuposClientID() error {
	if kc.config.ClientID == "" {
		return errors.New("keycloak client_id is empty")
	}

	params := url.Values{}
	params.Set("clientId", kc.config.ClientID)

	apiURL := fmt.Sprintf("%s/clients", kc.getAdminAPIURL())
	var clients []keycloakClientRepresentation
	if err := kc.doAdminGetRequest(apiURL, params, &clients); err != nil {
		return fmt.Errorf("query keycloak client by clientId %q: %w", kc.config.ClientID, err)
	}
	if len(clients) == 0 {
		return fmt.Errorf("keycloak client %q not found", kc.config.ClientID)
	}

	kc.config.SuposClientID = clients[0].ID
	return nil
}

func GetKeycloakClient() *KeycloakClient {
	if keycloakClient == nil {
		panic("Keycloak client not initialized")
	}
	return keycloakClient
}

// --- URL Helpers ---

func (kc *KeycloakClient) getAPIURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect", kc.config.IssuerURI, kc.config.Realm)
}

func (kc *KeycloakClient) getAdminAPIURL() string {
	return fmt.Sprintf("%s/admin/realms/%s", kc.config.IssuerURI, kc.config.Realm)
}

// --- Core API Methods ---

func (kc *KeycloakClient) GetAdminToken() (string, error) {
	if token, found := kc.tokenCache.Get(adminTokenCacheKey); found && token != nil {
		return token.(string), nil
	}

	apiURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", kc.config.IssuerURI)
	form := url.Values{
		"username":   {"admin"},
		"password":   {"tier0"},
		"grant_type": {"password"},
		"client_id":  {"admin-cli"},
	}

	var tokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, apiURL, form, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}

	kc.tokenCache.Set(adminTokenCacheKey, tokenResp.AccessToken, time.Duration(600)*time.Second)
	return tokenResp.AccessToken, nil
}

func (kc *KeycloakClient) Login(username, password string) (*AccessTokenDto, error) {
	form := url.Values{
		"grant_type":    {"password"},
		"username":      {username},
		"password":      {password},
		"client_id":     {kc.config.ClientID},
		"client_secret": {kc.config.ClientSecret},
	}
	var tokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/token", form, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (kc *KeycloakClient) GetKeyCloakTokenByCode(code string) (*AccessTokenDto, error) {
	form := url.Values{
		"grant_type":    {kc.config.AuthorizationGrantType},
		"code":          {code},
		"redirect_uri":  {kc.config.RedirectURI},
		"client_id":     {kc.config.ClientID},
		"client_secret": {kc.config.ClientSecret},
	}
	var tokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/token", form, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (kc *KeycloakClient) UserInfo(accessToken string) (*KeycloakUserInfoDto, error) {
	var userInfo KeycloakUserInfoDto
	if err := kc.doGetRequest(kc.getAPIURL()+"/userinfo", accessToken, nil, &userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (kc *KeycloakClient) RefreshToken(refreshToken string) (*AccessTokenDto, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {kc.config.ClientID},
		"client_secret": {kc.config.ClientSecret},
	}
	var tokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/token", form, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (kc *KeycloakClient) Logout(refreshToken string) error {
	form := url.Values{
		"client_id":     {kc.config.ClientID},
		"client_secret": {kc.config.ClientSecret},
		"refresh_token": {refreshToken},
	}
	return kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/logout", form, nil)
}

func (kc *KeycloakClient) CreateUser(user map[string]any) (string, error) {
	apiURL := kc.getAdminAPIURL() + "/users"
	resp, err := kc.doAdminJSONRequest(http.MethodPost, apiURL, user, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return "", berrors.NewBuzError(berrors.UserAlreadyExists, "user already exists")
	}
	if resp.StatusCode != http.StatusCreated {
		return "", handleAPIError("create user", resp)
	}

	location := resp.Header.Get("Location")
	parts := strings.Split(location, "/")
	return parts[len(parts)-1], nil
}

func (kc *KeycloakClient) FetchUser(username string) (*KeycloakUserInfoDto, error) {
	params := url.Values{"exact": {"true"}, "username": {username}}
	var users []KeycloakUserInfoDto
	if err := kc.doAdminGetRequest(kc.getAdminAPIURL()+"/users", params, &users); err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil // Not found
}

func (kc *KeycloakClient) FetchUserByEmail(email string) (*KeycloakUserInfoDto, error) {
	params := url.Values{"exact": {"true"}, "email": {email}}
	var users []KeycloakUserInfoDto
	if err := kc.doAdminGetRequest(kc.getAdminAPIURL()+"/users", params, &users); err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil // Not found
}

func (kc *KeycloakClient) DeleteUser(id string) error {
	apiURL := fmt.Sprintf("%s/users/%s", kc.getAdminAPIURL(), id)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

func (kc *KeycloakClient) ResetPassword(userID, password string) error {
	apiURL := fmt.Sprintf("%s/users/%s/reset-password", kc.getAdminAPIURL(), userID)
	body := map[string]any{
		"type":      "password",
		"temporary": false,
		"value":     password,
	}
	return kc.doAdminSimpleRequest(http.MethodPut, apiURL, body)
}

func (kc *KeycloakClient) UpdateUser(userID string, user map[string]any) error {
	apiURL := fmt.Sprintf("%s/users/%s", kc.getAdminAPIURL(), userID)
	return kc.doAdminSimpleRequest(http.MethodPut, apiURL, user)
}

func (kc *KeycloakClient) SetUserRoles(userID string, roles []KeycloakRoleInfoDto, add bool) error {
	apiURL := fmt.Sprintf("%s/users/%s/role-mappings/clients/%s", kc.getAdminAPIURL(), userID, kc.config.SuposClientID)
	method := http.MethodDelete
	if add {
		method = http.MethodPost
	}
	return kc.doAdminSimpleRequest(method, apiURL, roles)
}

func (kc *KeycloakClient) GetAllRoles() ([]KeycloakRoleInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/roles", kc.getAdminAPIURL(), kc.config.SuposClientID)
	var roles []KeycloakRoleInfoDto
	if err := kc.doAdminGetRequest(apiURL, nil, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (kc *KeycloakClient) FetchRole(roleName string) (*KeycloakRoleInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/roles/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, roleName)
	var role KeycloakRoleInfoDto
	err := kc.doAdminGetRequest(apiURL, nil, &role)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Response.StatusCode == http.StatusNotFound {
			return nil, nil // Not found, not an error
		}
		return nil, err
	}
	return &role, nil
}

func (kc *KeycloakClient) FetchRoleByID(id string) (*KeycloakRoleInfoDto, error) {
	apiURL := fmt.Sprintf("%s/roles-by-id/%s", kc.getAdminAPIURL(), id)
	var role KeycloakRoleInfoDto
	err := kc.doAdminGetRequest(apiURL, nil, &role)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Response.StatusCode == http.StatusNotFound {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &role, nil
}

func (kc *KeycloakClient) CreateRole(name, description string) error {
	apiURL := fmt.Sprintf("%s/clients/%s/roles", kc.getAdminAPIURL(), kc.config.SuposClientID)
	body := map[string]string{"name": name, "description": description}
	return kc.doAdminSimpleRequest(http.MethodPost, apiURL, body)
}

func (kc *KeycloakClient) DeleteRole(name string) error {
	apiURL := fmt.Sprintf("%s/clients/%s/roles/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, name)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

func (kc *KeycloakClient) GetRoleListByUserID(userID string) (map[string]any, error) {
	apiURL := fmt.Sprintf("%s/users/%s/role-mappings", kc.getAdminAPIURL(), userID)
	var roles map[string]any
	if err := kc.doAdminGetRequest(apiURL, nil, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (kc *KeycloakClient) GetUserExchangeTokenByID(userID string) (*AccessTokenDto, error) {
	// Step 1: Get a token with client_credentials grant type (service account token)
	// Note: In Java version, client_id and client_secret were hardcoded.
	// We use the configured ones for consistency, assuming the 'supos' client has service accounts enabled.
	serviceTokenForm := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {kc.config.ClientID},
		"client_secret": {kc.config.ClientSecret},
	}
	var serviceTokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/token", serviceTokenForm, &serviceTokenResp); err != nil {
		return nil, fmt.Errorf("failed to get service token for exchange: %w", err)
	}

	// Step 2: Exchange the service token for the user's token
	exchangeTokenForm := url.Values{
		"grant_type":        {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"client_id":         {kc.config.ClientID},
		"client_secret":     {kc.config.ClientSecret},
		"subject_token":     {serviceTokenResp.AccessToken},
		"requested_subject": {userID},
		"scope":             {"openid profile"},
	}
	var userTokenResp AccessTokenDto
	if err := kc.doFormRequest(http.MethodPost, kc.getAPIURL()+"/token", exchangeTokenForm, &userTokenResp); err != nil {
		return nil, fmt.Errorf("failed to exchange token for user %s: %w", userID, err)
	}
	return &userTokenResp, nil
}

func (kc *KeycloakClient) RemoveSession(sessionState string) error {
	apiURL := fmt.Sprintf("%s/sessions/%s", kc.getAdminAPIURL(), sessionState)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

// --- Authorization Services ---

func (kc *KeycloakClient) FetchResource(name string) (*KeycloakResourceInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/resource/search", kc.getAdminAPIURL(), kc.config.SuposClientID)
	params := url.Values{"name": {name}} // Java version had `exact:true`, but search endpoint doesn't support it. Name search is exact by default.

	var resource KeycloakResourceInfoDto
	if err := kc.doAdminGetRequest(apiURL, params, &resource); err != nil {
		return nil, err
	}
	if resource.ID == "" {
		return nil, nil // Not found
	}
	return &resource, nil
}

func (kc *KeycloakClient) CreateResource(name, resourceType string, uris []string) (*KeycloakResourceInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/resource", kc.getAdminAPIURL(), kc.config.SuposClientID)
	body := KeycloakResourceInfoDto{
		Name:        name,
		DisplayName: name,
		Type:        resourceType,
		URIs:        uris,
	}
	var createdResource KeycloakResourceInfoDto
	resp, err := kc.doAdminJSONRequest(http.MethodPost, apiURL, body, &createdResource)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, handleAPIError("create resource", resp)
	}

	return &createdResource, nil
}

func (kc *KeycloakClient) UpdateResource(id string, resource map[string]any) error {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/resource/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, id)
	return kc.doAdminSimpleRequest(http.MethodPut, apiURL, resource)
}

func (kc *KeycloakClient) DeleteResource(name string) error {
	resource, err := kc.FetchResource(name)
	if err != nil {
		return err
	}
	if resource == nil {
		return nil // Resource doesn't exist, which is a success for deletion.
	}

	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/resource/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, resource.ID)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

func (kc *KeycloakClient) FetchPolicy(name string) (*KeycloakPolicyInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/policy/search", kc.getAdminAPIURL(), kc.config.SuposClientID)
	params := url.Values{"name": {name}}

	var policy KeycloakPolicyInfoDto
	if err := kc.doAdminGetRequest(apiURL, params, &policy); err != nil {
		return nil, err
	}
	if policy.ID == "" {
		return nil, nil // Not found
	}
	return &policy, nil
}

func (kc *KeycloakClient) CreatePolicy(name, description, roleID string) (*KeycloakPolicyInfoDto, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/policy/role", kc.getAdminAPIURL(), kc.config.SuposClientID)

	// Construct the complex body required by the Keycloak API for a role-based policy
	body := map[string]any{
		"name":        name,
		"description": description,
		"type":        "role",
		"logic":       "POSITIVE",
		"roles": []map[string]any{
			{
				"id":       roleID,
				"required": false, // This field is present in the Java example
			},
		},
	}

	var createdPolicy KeycloakPolicyInfoDto
	resp, err := kc.doAdminJSONRequest(http.MethodPost, apiURL, body, &createdPolicy)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, handleAPIError("create policy", resp)
	}

	return &createdPolicy, nil
}

func (kc *KeycloakClient) DeletePolicy(name string) error {
	policy, err := kc.FetchPolicy(name)
	if err != nil {
		return err
	}
	if policy == nil {
		return nil // Policy doesn't exist, success for deletion.
	}

	// Note: The delete endpoint needs the policy's ID and its type in the URL.
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/policy/%s/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, policy.Type, policy.ID)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

func (kc *KeycloakClient) FetchPermission(name string) (map[string]any, error) {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/permission/search", kc.getAdminAPIURL(), kc.config.SuposClientID)
	params := url.Values{"name": {name}}

	var permission map[string]any
	if err := kc.doAdminGetRequest(apiURL, params, &permission); err != nil {
		return nil, err
	}
	if len(permission) == 0 {
		return nil, nil // Not found
	}
	return permission, nil
}

func (kc *KeycloakClient) CreatePermission(name, description, policyID, resourceID string) error {
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/permission/resource", kc.getAdminAPIURL(), kc.config.SuposClientID)

	body := map[string]any{
		"name":             name,
		"description":      description,
		"decisionStrategy": "UNANIMOUS",
		"policies":         []string{policyID},
		"resources":        []string{resourceID},
	}

	resp, err := kc.doAdminJSONRequest(http.MethodPost, apiURL, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleAPIError("create permission", resp)
	}
	return nil
}

func (kc *KeycloakClient) DeletePermission(name string) error {
	permission, err := kc.FetchPermission(name)
	if err != nil {
		return err
	}
	if permission == nil {
		return nil // Permission doesn't exist, success for deletion.
	}

	permissionID, ok := permission["id"].(string)
	if !ok {
		return fmt.Errorf("could not find id in fetched permission for %s", name)
	}

	// The Java version uses `/permission/resource/{id}` which seems to be the correct endpoint for resource-based permissions.
	apiURL := fmt.Sprintf("%s/clients/%s/authz/resource-server/permission/resource/%s", kc.getAdminAPIURL(), kc.config.SuposClientID, permissionID)
	return kc.doAdminSimpleRequest(http.MethodDelete, apiURL, nil)
}

// --- Realm Management & Utils ---

func (kc *KeycloakClient) SetLocale(locale string) error {
	apiURL := kc.getAdminAPIURL()

	// Step 1: Get current realm info
	var realmInfo map[string]any
	if err := kc.doAdminGetRequest(apiURL, nil, &realmInfo); err != nil {
		logx.Errorf("Failed to get realm info, skipping locale update.: %s", err)
		return err
	}

	// Step 2: Check if update is needed
	defaultLocale, _ := realmInfo["defaultLocale"].(string)
	supportedLocales, _ := realmInfo["supportedLocales"].([]any)
	var firstSupportedLocale string
	if len(supportedLocales) > 0 {
		firstSupportedLocale, _ = supportedLocales[0].(string)
	}

	if locale == defaultLocale || locale == firstSupportedLocale {
		logx.Debugf("Locale already set, no update needed.")
		return nil
	}

	// Step 3: Prepare and apply the update
	realmInfo["internationalizationEnabled"] = true
	realmInfo["defaultLocale"] = locale

	// Ensure the new locale is in the supported list
	localeExists := false
	for _, loc := range supportedLocales {
		if loc.(string) == locale {
			localeExists = true
			break
		}
	}
	if !localeExists {
		supportedLocales = append(supportedLocales, locale)
	}
	realmInfo["supportedLocales"] = supportedLocales

	if err := kc.doAdminSimpleRequest(http.MethodPut, apiURL, realmInfo); err != nil {
		logx.Errorf("Failed to set Keycloak locale.: %s", err)
		return err
	}

	logx.Debugf("Successfully set Keycloak locale.: %s", locale)
	return nil
}

// RemovePortIfDefault removes default http/https ports from a URL string.
func RemovePortIfDefault(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	port := parsedURL.Port()
	if (parsedURL.Scheme == "http" && port == "80") || (parsedURL.Scheme == "https" && port == "443") {
		parsedURL.Host = parsedURL.Hostname()
	}

	return parsedURL.String(), nil
}

// --- Request Helpers ---

func (kc *KeycloakClient) doFormRequest(method, url string, form url.Values, target any) error {
	req, err := http.NewRequest(method, url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return kc.doAndDecode(req, target)
}

func (kc *KeycloakClient) doGetRequest(url, token string, params url.Values, target any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return kc.doAndDecode(req, target)
}

func (kc *KeycloakClient) doAdminGetRequest(url string, params url.Values, target any) error {
	token, err := kc.GetAdminToken()
	if err != nil {
		return err
	}
	return kc.doGetRequest(url, token, params, target)
}

func (kc *KeycloakClient) doAdminJSONRequest(method, url string, body, target any) (*http.Response, error) {
	token, err := kc.GetAdminToken()
	if err != nil {
		return nil, err
	}
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if target != nil && resp.Body != nil {
		// We need to handle this carefully. The caller might need the raw response.
		// A better pattern is to decode *after* checking status.
		// For now, let's keep it simple. The caller must close the body if target is nil.
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil && err != io.EOF {
			// Don't close here, caller will
			return resp, err
		}
	}
	return resp, nil
}

func (kc *KeycloakClient) doAdminSimpleRequest(method, url string, body any) error {
	resp, err := kc.doAdminJSONRequest(method, url, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return handleAPIError("request", resp)
	}
	return nil
}

func (kc *KeycloakClient) doAndDecode(req *http.Request, target any) error {
	logx.Debugf("Making HTTP request: %s, %s", req.Method, req.URL.String())
	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return handleAPIError("request", resp)
	}
	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func handleAPIError(context string, resp *http.Response) error {
	bodyBytes, _ := io.ReadAll(resp.Body)
	logx.Errorf("Keycloak API error: %s, %d, %s", context, resp.StatusCode, string(bodyBytes))
	return &APIError{Response: resp, Body: bodyBytes}
}
