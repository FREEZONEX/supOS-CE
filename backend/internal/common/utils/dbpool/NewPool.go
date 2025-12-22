package dbpool

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

var poolMap = make(map[string]*pgxpool.Pool, 4)
var poolLock sync.RWMutex
var initOnce sync.Once

func NewPool(ctx context.Context, connString, appName string) (*pgxpool.Pool, error) {
	poolLock.Lock()
	defer poolLock.Unlock()
	pool, err := newPool(ctx, connString, appName)
	if err != nil {
		return nil, err
	}
	poolMap[appName] = pool
	initOnce.Do(func() {
		go statsPool()
	})
	return pool, nil
}
func statsPool() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			poolLock.RLock()
			for name, pool := range poolMap {
				stats := pool.Stat()
				logx.Infof("[%s Stats] Total:%d Idle:%d Acquired:%d Max:%d", name,
					stats.TotalConns(), stats.IdleConns(),
					stats.AcquiredConns(), stats.MaxConns())
			}
			poolLock.RUnlock()
		}
	}
}
func newPool(ctx context.Context, connString, appName string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MinConns = 3 // 保持最小热连接
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
