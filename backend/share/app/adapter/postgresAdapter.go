package adapter

import (
	"backend/share/app/model"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL驱动
)

// PostgresConfig PostgreSQL连接配置
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	SSLMode  string
}

// PostgresAdapter PostgreSQL适配器
type PostgresAdapter struct {
	config *PostgresConfig
	db     *sql.DB
}

// NewPostgresAdapter 创建PostgreSQL适配器
func NewPostgresAdapter(config *PostgresConfig) *PostgresAdapter {
	return &PostgresAdapter{
		config: config,
	}
}

// DefaultPostgresConfig 默认PostgreSQL配置
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     "postgresql",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		SSLMode:  "disable",
	}
}

// Connect 连接到PostgreSQL数据库
func (pa *PostgresAdapter) Connect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		pa.config.Host, pa.config.Port, pa.config.User, pa.config.Password, pa.config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	pa.db = db
	log.Printf("成功连接到PostgreSQL: %s:%d", pa.config.Host, pa.config.Port)
	return nil
}

// Close 关闭数据库连接
func (pa *PostgresAdapter) Close() error {
	if pa.db != nil {
		return pa.db.Close()
	}
	return nil
}

// CreateDatabasesFromRequirements 根据requirement配置创建数据库
func (pa *PostgresAdapter) CreateDatabasesFromRequirements(requirements *model.Requirements) ([]*model.DatabaseInfo, error) {
	if requirements == nil || !requirements.HasDatabases() {
		return nil, fmt.Errorf("requirements配置为空")
	}

	log.Printf("开始处理 %d 个数据库需求", len(requirements.Databases))
	var dbInfoList []*model.DatabaseInfo
	// 遍历所有数据库需求
	for i, dbReq := range requirements.Databases {
		log.Printf("处理数据库需求 [%d/%d]: %s", i+1, len(requirements.Databases), dbReq.Name)
		dbInfo, err := pa.createDatabaseWithUser(dbReq)
		if err != nil {
			log.Printf("创建数据库 %s 失败: %v", dbReq.Name, err)
			return nil, fmt.Errorf("创建数据库 %s 失败: %v", dbReq.Name, err)
		}
		dbInfoList = append(dbInfoList, dbInfo)
	}

	log.Println("所有数据库需求处理完成")
	return dbInfoList, nil
}

// createDatabaseWithUser 创建数据库和对应的用户
func (pa *PostgresAdapter) createDatabaseWithUser(dbReq model.DatabaseRequirement) (*model.DatabaseInfo, error) {
	// 生成用户名和密码
	username := generateUsername(dbReq.Name)
	var password string

	// 1. 检查数据库是否已存在
	exists, err := pa.databaseExists(dbReq.Name)
	if err == nil && !exists {
		// 创建数据库
		if err := pa.createDatabase(dbReq.Name, username); err != nil {
			return nil, fmt.Errorf("创建数据库失败: %v", err)
		}
		log.Printf("创建数据库: %s", dbReq.Name)
	}

	// 2. 检查用户是否已存在
	userExists, err := pa.userExists(username)
	if err != nil {
		return nil, fmt.Errorf("检查用户是否存在失败: %v", err)
	}

	// 3. 创建用户（如果不存在）
	if !userExists {
		password = generatePassword()
		if err := pa.createUser(username, password); err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
		log.Printf("创建用户: %s", username)
	} else {
		log.Printf("用户 %s 已存在，使用现有用户", username)
	}

	// 4. 授予权限
	if err := pa.grantPrivileges(dbReq.Name, username); err != nil {
		return nil, fmt.Errorf("授予权限失败: %v", err)
	}
	log.Printf("为用户 %s 授予数据库 %s 的权限", username, dbReq.Name)

	// 5. 记录数据库信息
	dbInfo := model.NewPostgresDatabaseInfo(dbReq.Name, username, password, pa.config.Host, fmt.Sprintf("%d", pa.config.Port))

	log.Printf("数据库创建完成: %v", dbInfo)

	return dbInfo, nil
}

// databaseExists 检查数据库是否存在
func (pa *PostgresAdapter) databaseExists(dbName string) (bool, error) {
	query := "SELECT 1 FROM pg_database WHERE datname = $1"
	var exists int
	err := pa.db.QueryRow(query, dbName).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// userExists 检查用户是否存在
func (pa *PostgresAdapter) userExists(username string) (bool, error) {
	query := "SELECT 1 FROM pg_roles WHERE rolname = $1"
	var exists int
	err := pa.db.QueryRow(query, username).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// createUser 创建用户
func (pa *PostgresAdapter) createUser(username, password string) error {
	// 使用参数化查询避免SQL注入
	query := fmt.Sprintf("CREATE USER %s WITH PASSWORD $1", quoteIdentifier(username))
	_, err := pa.db.Exec(query, password)
	return err
}

// createDatabase 创建数据库
func (pa *PostgresAdapter) createDatabase(dbName, owner string) error {
	// 创建数据库并指定所有者
	query := fmt.Sprintf("CREATE DATABASE %s OWNER %s",
		quoteIdentifier(dbName), quoteIdentifier(owner))
	_, err := pa.db.Exec(query)
	return err
}

// grantPrivileges 授予权限
func (pa *PostgresAdapter) grantPrivileges(dbName, username string) error {
	// 授予所有权限
	queries := []string{
		fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s",
			quoteIdentifier(dbName), quoteIdentifier(username)),
		fmt.Sprintf("ALTER DATABASE %s OWNER TO %s",
			quoteIdentifier(dbName), quoteIdentifier(username)),
	}

	for _, query := range queries {
		if _, err := pa.db.Exec(query); err != nil {
			return fmt.Errorf("执行SQL失败: %s, 错误: %v", query, err)
		}
	}

	return nil
}

// quoteIdentifier 引用标识符（防止SQL注入）
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// generateUsername 生成用户名
func generateUsername(dbName string) string {
	// 将数据库名转换为有效的用户名
	username := strings.ToLower(dbName)
	username = strings.ReplaceAll(username, "-", "_")
	username = strings.ReplaceAll(username, ".", "_")

	// 确保用户名以字母开头
	if len(username) > 0 && !isLetter(rune(username[0])) {
		username = "u_" + username
	}

	// 限制长度
	if len(username) > 32 {
		username = username[:32]
	}

	return username
}

// generatePassword 生成随机密码
func generatePassword() string {
	// 简单的密码生成逻辑，实际使用时应该使用更安全的随机生成器
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("Pwd_%d_@1304", timestamp)
}

// isLetter 检查字符是否为字母
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// TestConnection 测试数据库连接
func (pa *PostgresAdapter) TestConnection() error {
	if pa.db == nil {
		return fmt.Errorf("数据库未连接")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pa.db.PingContext(ctx)
}

// GetDatabaseInfo 获取数据库信息
func (pa *PostgresAdapter) GetDatabaseInfo(dbName string) (map[string]string, error) {
	query := `
		SELECT 
			d.datname as database,
			u.usename as owner,
			d.datcollate as collation,
			d.datctype as ctype,
			d.datconnlimit as conn_limit,
			d.datallowconn as allow_connections
		FROM pg_database d
		JOIN pg_user u ON d.datdba = u.usesysid
		WHERE d.datname = $1
	`

	var (
		database, owner, collation, ctype string
		connLimit                         int
		allowConnections                  bool
	)

	err := pa.db.QueryRow(query, dbName).Scan(
		&database, &owner, &collation, &ctype, &connLimit, &allowConnections,
	)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"database":          database,
		"owner":             owner,
		"collation":         collation,
		"ctype":             ctype,
		"connection_limit":  fmt.Sprintf("%d", connLimit),
		"allow_connections": fmt.Sprintf("%v", allowConnections),
	}, nil
}

// ListDatabases 列出所有数据库
func (pa *PostgresAdapter) ListDatabases() ([]string, error) {
	query := "SELECT datname FROM pg_database WHERE datistemplate = false"
	rows, err := pa.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

// CleanupTestDatabases 清理测试数据库（用于测试）
func (pa *PostgresAdapter) CleanupTestDatabases(prefix string) error {
	databases, err := pa.ListDatabases()
	if err != nil {
		return err
	}

	for _, dbName := range databases {
		if strings.HasPrefix(dbName, prefix) {
			log.Printf("清理测试数据库: %s", dbName)
			if err := pa.dropDatabase(dbName); err != nil {
				log.Printf("删除数据库 %s 失败: %v", dbName, err)
			}
		}
	}

	return nil
}

// dropDatabase 删除数据库
func (pa *PostgresAdapter) dropDatabase(dbName string) error {
	// 首先断开所有连接到该数据库的连接
	disconnectQuery := `
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = $1 AND pid <> pg_backend_pid()
	`
	_, _ = pa.db.Exec(disconnectQuery, dbName)

	// 删除数据库
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", quoteIdentifier(dbName))
	_, err := pa.db.Exec(query)
	return err
}
