package postgresql

import (
	"backend/internal/types"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestListTableInfo(t *testing.T) {
	// 初始化数据库连接池
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.22:34099/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// 示例数据
	topics := []*types.CreateTopicDto{
		{TableName: "supos.uns_namespace"},
		{TableName: "supos.supos_example"},
	}

	// 查询表信息
	tableInfos, err := ListTableInfos(pool, topics)
	if err != nil {
		panic(err)
	}

	// 打印结果
	for tableName, info := range tableInfos {
		fmt.Printf("Table: %s\n", tableName)
		fmt.Printf("Primary Keys: %v\n", info.PKs)

		for key, value := range info.FieldTypes {
			fmt.Printf("  %s: %s\n", key, value)
		}
		fmt.Println()
	}
}
func TestPgTempTable(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	defer conn.Release()
	sql := "CREATE TEMP TABLE temp__f1_c9 (LIKE public.uns_label_ref EXCLUDING INDEXES) "
	tag, err := conn.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}
	utcTime := time.Now().UTC()
	ts := utcTime.Format("2006-01-02 15:04:05.000") + "+00"
	t.Log("ts:", ts)
	tag, err = conn.Exec(context.Background(), `insert into temp__f1_c9("label_id","uns_id") values (200,2),(200,2),(211,33)`)
	if err != nil {
		panic(err)
	}
	t.Log(tag.RowsAffected())

	tag, err = conn.Exec(context.Background(), `insert into public.uns_label_ref("label_id","uns_id") 
    select "label_id","uns_id" from temp__f1_c9 ON CONFLICT("label_id","uns_id") DO NOTHING `)
	if err != nil {
		panic(err)
	}
	t.Log(tag.RowsAffected())
}
