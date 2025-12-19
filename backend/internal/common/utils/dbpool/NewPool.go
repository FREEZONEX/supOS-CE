package dbpool

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connString, appName string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MinConns = 3 // 保持最小热连接
	if _, has := config.ConnConfig.RuntimeParams["statement_timeout"]; !has {
		config.ConnConfig.RuntimeParams["statement_timeout"] = "10000" // 10秒
	}
	// 1. 设置应用名，便于在数据库端监控和区分连接来源
	config.ConnConfig.RuntimeParams["application_name"] = appName

	if config.ConnConfig.ConnectTimeout == 0 {
		config.ConnConfig.ConnectTimeout = time.Second * 5
	}
	if config.MaxConnIdleTime == 0 || config.MaxConnIdleTime == time.Minute*30 {
		config.MaxConnIdleTime = 15 * time.Minute
	}
	if config.MaxConnLifetime == time.Hour {
		config.MaxConnLifetime = time.Minute * 30
	}
	if config.HealthCheckPeriod == time.Minute {
		config.HealthCheckPeriod = 15 * time.Second
	}
	//2.  初始化设置
	config.ConnConfig.AfterConnect = func(ctx context.Context, conn *pgconn.PgConn) (checkErr error) {
		// 新连接建立后的初始化，例如设置时区、搜索路径
		// 可以包含一个简单查询，但这只是初始验证
		checkErr = conn.Exec(ctx, "SET timezone = 'UTC';").Close()
		return checkErr
	}
	//3.  获取连接前的最后检查
	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// 在实际使用连接前做最终校验
		var one int
		err := conn.QueryRow(ctx, "SELECT 1").Scan(&one)
		return err == nil // 返回 false 则丢弃此连接
	}

	return pgxpool.NewWithConfig(ctx, config)
}
