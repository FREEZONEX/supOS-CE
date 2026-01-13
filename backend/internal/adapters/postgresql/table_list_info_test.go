package postgresql

import (
	"backend/internal/common"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestListTableInfo(t *testing.T) {
	// 初始化数据库连接池
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	{
		_, err = pool.Exec(context.Background(), `insert into no121("timeStamp","json","id") values (NOW(),'{}',1)`)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				t.Logf("no121 pgErr: code=%v, msg=%v", pgErr.Code, pgErr.Message)
			}
		}
	}
	{
		_, err = pool.Exec(context.Background(), `insert into _121_77ffd11e5c784f3b8396("timeStamp","json","id") values (NOW(),'{}',1)`)
		if err != nil {
			if pgErr, is := err.(*pgconn.PgError); is {
				t.Logf("noid: pgErr: code=%v, msg=%v", pgErr.Code, pgErr.Message)
			} else {
				t.Error("noid:", err)
			}
		}
	}
	{
		_, err = pool.Exec(context.Background(), `insert into _121_77ffd11e5c784f3b8396("timeStamp","json","_id") values (NOW(),$1,$2)`, `{}`)
		if err != nil {
			if pgErr, is := err.(*pgconn.PgError); is {
				t.Logf("mismatched: pgErr: code=%v, msg=%v", pgErr.Code, pgErr.Message)
			} else {
				t.Error("mismatched:", err)
			}
		}
	}
	// 示例数据
	// 查询表信息
	tableInfos, err := ListTableInfos(pool, []string{"supos.uns_namespace", "supos.supos_example"})
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
func TestFillLastRecord(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@100.100.100.20:31014/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	fields := []*types.FieldDefine{
		{
			Name:   "timeStamp",
			Type:   types.FieldTypeDatetime,
			Unique: base.OptionalTrue,
		},
		{
			Name:   "tm",
			Type:   types.FieldTypeLong,
			Unique: base.OptionalTrue,
		},
		{
			Name: "wq",
			Type: types.FieldTypeDouble,
		},
		{
			Name: "status",
			Type: types.FieldTypeLong,
		}}
	uns := &types.CreateTopicDto{
		Alias:  "akseqxajk8",
		Fields: fields,
	}
	common.InitSnowflake(123)
	pgQuery := func(ctx context.Context) (pgx.Rows, error) {
		return pool.Query(ctx, fmt.Sprintf(`select * from "%s" ORDER BY "%s" DESC LIMIT 1`, uns.GetTable(), uns.GetTimestampField()))
	}
	FillLastRecord(logx.WithContext(context.Background()), uns, pgQuery)
}
