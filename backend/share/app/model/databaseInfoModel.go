package model

// DatabaseInfo 表示数据库连接信息
type DatabaseInfo struct {
	// 数据库名称
	Database string `json:"database" yaml:"database" ini:"database"`
	// 用户名
	Username string `json:"username" yaml:"username" ini:"username"`
	// 密码
	Password string `json:"password" yaml:"password" ini:"password"`
	// 主机地址
	Host string `json:"host" yaml:"host" ini:"host"`
	// 端口号
	Port string `json:"port" yaml:"port" ini:"port"`
	// 连接字符串（可选，可自动生成）
	ConnectionString string `json:"connectionString,omitempty" yaml:"connectionString,omitempty" ini:"connectionString,omitempty"`
	// 数据库类型（如postgresql、mysql等）
	DatabaseType string `json:"databaseType,omitempty" yaml:"databaseType,omitempty" ini:"databaseType,omitempty"`
	// SSL模式（仅PostgreSQL）
	SSLMode string `json:"sslMode,omitempty" yaml:"sslMode,omitempty" ini:"sslMode,omitempty"`
	// 字符集
	Charset string `json:"charset,omitempty" yaml:"charset,omitempty" ini:"charset,omitempty"`
	// 时区
	Timezone string `json:"timezone,omitempty" yaml:"timezone,omitempty" ini:"timezone,omitempty"`
}

// NewDatabaseInfo 创建新的数据库信息
func NewDatabaseInfo(database, username, password, host, port string) *DatabaseInfo {
	return &DatabaseInfo{
		Database: database,
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

// NewPostgresDatabaseInfo 创建PostgreSQL数据库信息
func NewPostgresDatabaseInfo(database, username, password, host, port string) *DatabaseInfo {
	dbInfo := NewDatabaseInfo(database, username, password, host, port)
	dbInfo.DatabaseType = "postgresql"
	dbInfo.SSLMode = "disable" // 默认禁用SSL
	dbInfo.ConnectionString = dbInfo.GeneratePostgresConnectionString()
	return dbInfo
}

// GeneratePostgresConnectionString 生成PostgreSQL连接字符串
func (d *DatabaseInfo) GeneratePostgresConnectionString() string {
	if d.DatabaseType == "postgresql" {
		return d.generatePostgresDSN()
	}
	return d.generateGenericDSN()
}

// generatePostgresDSN 生成PostgreSQL DSN
func (d *DatabaseInfo) generatePostgresDSN() string {
	dsn := "postgresql://"
	if d.Username != "" && d.Password != "" {
		dsn += d.Username + ":" + d.Password + "@"
	}
	dsn += d.Host + ":" + d.Port + "/" + d.Database

	params := []string{}
	if d.SSLMode != "" {
		params = append(params, "sslmode="+d.SSLMode)
	}
	if d.Charset != "" {
		params = append(params, "charset="+d.Charset)
	}
	if d.Timezone != "" {
		params = append(params, "timezone="+d.Timezone)
	}

	if len(params) > 0 {
		dsn += "?" + joinParams(params, "&")
	}

	return dsn
}

// generateGenericDSN 生成通用DSN
func (d *DatabaseInfo) generateGenericDSN() string {
	dsn := d.DatabaseType + "://"
	if d.Username != "" && d.Password != "" {
		dsn += d.Username + ":" + d.Password + "@"
	}
	dsn += d.Host + ":" + d.Port + "/" + d.Database

	params := []string{}
	if d.Charset != "" {
		params = append(params, "charset="+d.Charset)
	}
	if d.Timezone != "" {
		params = append(params, "timezone="+d.Timezone)
	}

	if len(params) > 0 {
		dsn += "?" + joinParams(params, "&")
	}

	return dsn
}

// ToMap 将DatabaseInfo转换为map[string]string
func (d *DatabaseInfo) ToMap() map[string]string {
	return map[string]string{
		"database":         d.Database,
		"username":         d.Username,
		"password":         d.Password,
		"host":             d.Host,
		"port":             d.Port,
		"databaseType":     d.DatabaseType,
		"sslMode":          d.SSLMode,
		"charset":          d.Charset,
		"timezone":         d.Timezone,
		"connectionString": d.ConnectionString,
	}
}

// FromMap 从map[string]string创建DatabaseInfo
func FromMap(data map[string]string) *DatabaseInfo {
	dbInfo := &DatabaseInfo{
		Database: data["database"],
		Username: data["username"],
		Password: data["password"],
		Host:     data["host"],
		Port:     data["port"],
	}

	// 可选字段
	if dbType, ok := data["databaseType"]; ok {
		dbInfo.DatabaseType = dbType
	}
	if sslMode, ok := data["sslMode"]; ok {
		dbInfo.SSLMode = sslMode
	}
	if charset, ok := data["charset"]; ok {
		dbInfo.Charset = charset
	}
	if timezone, ok := data["timezone"]; ok {
		dbInfo.Timezone = timezone
	}
	if connStr, ok := data["connectionString"]; ok {
		dbInfo.ConnectionString = connStr
	} else {
		dbInfo.ConnectionString = dbInfo.GeneratePostgresConnectionString()
	}

	return dbInfo
}

// Validate 验证数据库信息
func (d *DatabaseInfo) Validate() error {
	if d.Database == "" {
		return NewValidationError("database name cannot be empty")
	}
	if d.Username == "" {
		return NewValidationError("username cannot be empty")
	}
	if d.Password == "" {
		return NewValidationError("password cannot be empty")
	}
	if d.Host == "" {
		return NewValidationError("host cannot be empty")
	}
	if d.Port == "" {
		return NewValidationError("port cannot be empty")
	}
	return nil
}

// GetConnectionInfo 获取连接信息（隐藏密码）
func (d *DatabaseInfo) GetConnectionInfo() map[string]string {
	return map[string]string{
		"database":     d.Database,
		"username":     d.Username,
		"host":         d.Host,
		"port":         d.Port,
		"databaseType": d.DatabaseType,
		"sslMode":      d.SSLMode,
		"charset":      d.Charset,
		"timezone":     d.Timezone,
	}
}

// MaskPassword 返回一个隐藏密码的副本
func (d *DatabaseInfo) MaskPassword() *DatabaseInfo {
	return &DatabaseInfo{
		Database:         d.Database,
		Username:         d.Username,
		Password:         "******",
		Host:             d.Host,
		Port:             d.Port,
		DatabaseType:     d.DatabaseType,
		SSLMode:          d.SSLMode,
		Charset:          d.Charset,
		Timezone:         d.Timezone,
		ConnectionString: d.ConnectionString,
	}
}

// DatabaseInfoList 数据库信息列表
type DatabaseInfoList struct {
	Databases []*DatabaseInfo `json:"databases" yaml:"databases"`
}

// NewDatabaseInfoList 创建新的数据库信息列表
func NewDatabaseInfoList() *DatabaseInfoList {
	return &DatabaseInfoList{
		Databases: make([]*DatabaseInfo, 0),
	}
}

// Add 添加数据库信息
func (l *DatabaseInfoList) Add(dbInfo *DatabaseInfo) {
	l.Databases = append(l.Databases, dbInfo)
}

// FindByDatabaseName 根据数据库名称查找
func (l *DatabaseInfoList) FindByDatabaseName(name string) *DatabaseInfo {
	for _, dbInfo := range l.Databases {
		if dbInfo.Database == name {
			return dbInfo
		}
	}
	return nil
}

// FindByUsername 根据用户名查找
func (l *DatabaseInfoList) FindByUsername(username string) []*DatabaseInfo {
	var result []*DatabaseInfo
	for _, dbInfo := range l.Databases {
		if dbInfo.Username == username {
			result = append(result, dbInfo)
		}
	}
	return result
}

// RemoveByDatabaseName 根据数据库名称移除
func (l *DatabaseInfoList) RemoveByDatabaseName(name string) bool {
	for i, dbInfo := range l.Databases {
		if dbInfo.Database == name {
			l.Databases = append(l.Databases[:i], l.Databases[i+1:]...)
			return true
		}
	}
	return false
}

// ToMapSlice 转换为map切片
func (l *DatabaseInfoList) ToMapSlice() []map[string]string {
	var result []map[string]string
	for _, dbInfo := range l.Databases {
		result = append(result, dbInfo.ToMap())
	}
	return result
}

// helper function to join parameters
func joinParams(params []string, sep string) string {
	if len(params) == 0 {
		return ""
	}
	if len(params) == 1 {
		return params[0]
	}

	result := params[0]
	for i := 1; i < len(params); i++ {
		result += sep + params[i]
	}
	return result
}
