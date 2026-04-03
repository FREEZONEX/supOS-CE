package relationDB

import "time"

// IamUser stores core IAM user data.
type IamUser struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64)"`
	Username       string    `gorm:"column:username;type:varchar(200);not null;uniqueIndex:idx_supos_user_username"`
	DisplayName    string    `gorm:"column:display_name;type:varchar(200)"`
	Email          string    `gorm:"column:email;type:varchar(255)"`
	Enabled        bool      `gorm:"column:enabled;not null;default:true"`
	Source         string    `gorm:"column:source;type:varchar(64)"`
	Password       string    `gorm:"column:password;type:text"`
	Phone          string    `gorm:"column:phone;type:varchar(64)"`
	HomePage       string    `gorm:"column:home_page;type:varchar(255)"`
	FirstTimeLogin int       `gorm:"column:first_time_login;not null;default:1"`
	TipsEnable     int       `gorm:"column:tips_enable;not null;default:1"`
	MainLanguage   string    `gorm:"column:main_language;type:varchar(32)"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (IamUser) TableName() string {
	return "supos_user"
}

// IamRole stores IAM role definitions.
type IamRole struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(64)"`
	RoleKey     string    `gorm:"column:role_key;type:varchar(64);not null;uniqueIndex:idx_supos_role_key"`
	RoleName    string    `gorm:"column:role_name;type:varchar(200);not null"`
	Description string    `gorm:"column:description;type:text"`
	Builtin     bool      `gorm:"column:builtin;not null;default:false"`
	Status      int       `gorm:"column:status;not null;default:1"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (IamRole) TableName() string {
	return "supos_role"
}

// IamUserRole stores user-role assignments.
type IamUserRole struct {
	UserID    string    `gorm:"column:user_id;primaryKey;type:varchar(64)"`
	RoleID    string    `gorm:"column:role_id;primaryKey;type:varchar(64)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy string    `gorm:"column:created_by;type:varchar(64)"`
}

func (IamUserRole) TableName() string {
	return "supos_user_role"
}

// IamRoleResource stores role-resource assignments.
type IamRoleResource struct {
	RoleID     string    `gorm:"column:role_id;primaryKey;type:varchar(64)"`
	ResourceID int64     `gorm:"column:resource_id;primaryKey"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy  string    `gorm:"column:created_by;type:varchar(64)"`
}

func (IamRoleResource) TableName() string {
	return "supos_role_resource"
}

// IamSession stores platform-owned session state.
type IamSession struct {
	ID            string     `gorm:"column:id;primaryKey;type:varchar(128)"`
	UserID        string     `gorm:"column:user_id;type:varchar(64);not null;index:idx_supos_session_user"`
	ExpiredAt     time.Time  `gorm:"column:expired_at;not null;index:idx_supos_session_expired"`
	LastAccessAt  time.Time  `gorm:"column:last_access_at;not null"`
	RevokedAt     *time.Time `gorm:"column:revoked_at"`
	LegacyToken   string     `gorm:"column:legacy_token;type:text"`
	PasswordReady bool       `gorm:"column:password_ready;not null;default:false"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (IamSession) TableName() string {
	return "supos_session"
}

// IamOAuthClient stores minimal OAuth client configuration for internal integrations.
type IamOAuthClient struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ClientID     string    `gorm:"column:client_id;type:varchar(128);not null;uniqueIndex:idx_supos_oauth_client_id"`
	ClientSecret string    `gorm:"column:client_secret;type:text"`
	ClientName   string    `gorm:"column:client_name;type:varchar(200)"`
	RedirectURIs string    `gorm:"column:redirect_uris;type:text"`
	Scopes       string    `gorm:"column:scopes;type:text"`
	GrantTypes   string    `gorm:"column:grant_types;type:text"`
	Enabled      bool      `gorm:"column:enabled;not null;default:true"`
	Trusted      bool      `gorm:"column:trusted;not null;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (IamOAuthClient) TableName() string {
	return "supos_oauth_client"
}

// IamOAuthAuthorizationCode stores short-lived authorization codes.
type IamOAuthAuthorizationCode struct {
	Code        string     `gorm:"column:code;primaryKey;type:varchar(256)"`
	ClientID    string     `gorm:"column:client_id;type:varchar(128);not null;index:idx_supos_oauth_code_client"`
	UserID      string     `gorm:"column:user_id;type:varchar(64);not null"`
	RedirectURI string     `gorm:"column:redirect_uri;type:text"`
	Scopes      string     `gorm:"column:scopes;type:text"`
	ExpiredAt   time.Time  `gorm:"column:expired_at;not null;index:idx_supos_oauth_code_expired"`
	UsedAt      *time.Time `gorm:"column:used_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (IamOAuthAuthorizationCode) TableName() string {
	return "supos_oauth_authorization_code"
}

// IamOAuthAccessToken stores opaque access tokens for internal OAuth consumers.
type IamOAuthAccessToken struct {
	AccessToken string     `gorm:"column:access_token;primaryKey;type:varchar(256)"`
	ClientID    string     `gorm:"column:client_id;type:varchar(128);not null;index:idx_supos_oauth_token_client"`
	UserID      string     `gorm:"column:user_id;type:varchar(64);not null"`
	Scopes      string     `gorm:"column:scopes;type:text"`
	ExpiredAt   time.Time  `gorm:"column:expired_at;not null;index:idx_supos_oauth_token_expired"`
	RevokedAt   *time.Time `gorm:"column:revoked_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (IamOAuthAccessToken) TableName() string {
	return "supos_oauth_access_token"
}
